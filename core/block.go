package core

import (
	"errors"
	"log"
	"strconv"

	"github.com/nikhan/go-fetch"
)

// NewBlock creates a new block from a spec
func NewBlock(s Spec) *Block {
	var in []Input
	var out []Output

	for _, v := range s.Inputs {

		q, _ := fetch.Parse(".")
		in = append(in, Input{
			Name:  v.Name,
			Value: q,
			C:     make(chan Message),
		})
	}

	for _, v := range s.Outputs {
		out = append(out, Output{
			Name:        v.Name,
			Connections: make(map[Connection]struct{}),
		})
	}

	return &Block{
		state: BlockState{
			make(MessageMap),
			make(MessageMap),
			make(MessageMap),
			make(Manifest),
			false,
		},
		routing: BlockRouting{
			Inputs:        in,
			Outputs:       out,
			InterruptChan: make(chan Interrupt),
		},
		kernel:     s.Kernel,
		sourceType: s.Source,
	}
}

// suture: the main routine the block runs
func (b *Block) Serve() {
	for {
		var interrupt Interrupt

		b.routing.RLock()
		for {
			interrupt = b.receive()
			if interrupt != nil {
				break
			}

			interrupt = b.process()
			if interrupt != nil {
				break
			}

			interrupt = b.broadcast()
			if interrupt != nil {
				break
			}

			b.crank()
		}
		b.routing.RUnlock()
		b.routing.Lock()
		if ok := interrupt(); !ok {
			return
		}
		b.routing.Unlock()
	}
}

func (b *Block) exportInput(id RouteIndex) (*Input, error) {
	if int(id) >= len(b.routing.Inputs) || int(id) < 0 {
		return nil, errors.New("index out of range")
	}

	var v interface{}
	switch n := b.routing.Inputs[id].Value.(type) {
	case *fetch.Query:
		// yuck copy
		v, _ = fetch.Parse(n.String())
	default:
		v = Copy(n)
	}

	return &Input{
		Value: v,
		C:     b.routing.Inputs[id].C,
		Name:  b.routing.Inputs[id].Name,
	}, nil
}

// GetInput returns the specified Input
func (b *Block) GetInput(id RouteIndex) (Input, error) {
	b.routing.RLock()
	r, err := b.exportInput(id)
	b.routing.RUnlock()
	return *r, err
}

// GetInputs returns all inputs for a block.
func (b *Block) GetInputs() []Input {
	b.routing.RLock()
	re := make([]Input, len(b.routing.Inputs), len(b.routing.Inputs))
	for i, _ := range b.routing.Inputs {
		r, _ := b.exportInput(RouteIndex(i))
		re[i] = *r
	}
	b.routing.RUnlock()
	return re
}

// Outputs return a list of manifest pairs for the block
func (b *Block) GetOutputs() []Output {
	b.routing.RLock()
	m := make([]Output, len(b.routing.Outputs), len(b.routing.Outputs))
	for id, out := range b.routing.Outputs {
		m[id] = Output{
			Name:        out.Name,
			Connections: make(map[Connection]struct{}),
		}
		for k, _ := range out.Connections {
			m[id].Connections[k] = struct{}{}
		}
	}
	b.routing.RUnlock()
	return m
}

func (b *Block) GetSource() Source {
	b.routing.RLock()
	v := b.routing.Source
	b.routing.RUnlock()
	return v
}

// sets a store for the block. can be set to nil
func (b *Block) SetSource(s Source) error {
	returnVal := make(chan error, 1)
	b.routing.InterruptChan <- func() bool {
		if s != nil && s.GetType() != b.sourceType {
			log.Println(s.GetType())
			returnVal <- errors.New("invalid source type (" + strconv.Itoa(int(s.GetType())) + ") for this block (needs " + strconv.Itoa(int(b.sourceType)) + ")")
			return true
		}
		b.routing.Source = s
		returnVal <- nil
		return true
	}
	return <-returnVal
}

// RouteValue sets the route to always be the specified value
func (b *Block) SetInput(id RouteIndex, v interface{}) error {
	returnVal := make(chan error, 1)
	b.routing.InterruptChan <- func() bool {
		if int(id) < 0 || int(id) >= len(b.routing.Inputs) {
			returnVal <- errors.New("input out of range")
			return true
		}

		b.routing.Inputs[id].Value = v
		returnVal <- nil
		return true
	}
	return <-returnVal
}

// Connect connects a Route, specified by ID, to a connection
func (b *Block) Connect(id RouteIndex, c Connection) error {
	returnVal := make(chan error, 1)
	b.routing.InterruptChan <- func() bool {
		if int(id) < 0 || int(id) >= len(b.routing.Outputs) {
			returnVal <- errors.New("output out of range")
			return true
		}

		if _, ok := b.routing.Outputs[id].Connections[c]; ok {
			returnVal <- errors.New("this connection already exists on this output")
			return true
		}

		b.routing.Outputs[id].Connections[c] = struct{}{}
		returnVal <- nil
		return true
	}
	return <-returnVal
}

// Disconnect removes a connection from a Input
func (b *Block) Disconnect(id RouteIndex, c Connection) error {
	returnVal := make(chan error, 1)
	b.routing.InterruptChan <- func() bool {
		if int(id) < 0 || int(id) >= len(b.routing.Outputs) {
			returnVal <- errors.New("output out of range")
			return true
		}

		if _, ok := b.routing.Outputs[id].Connections[c]; !ok {
			returnVal <- errors.New("connection does not exist")
			return true
		}

		delete(b.routing.Outputs[id].Connections, c)
		returnVal <- nil
		return true
	}
	return <-returnVal
}

// suture: stop the block
func (b *Block) Stop() {
	b.routing.InterruptChan <- func() bool {
		return false
	}
}

// wait and listen for all kernel inputs to be filled.
func (b *Block) receive() Interrupt {
	var err error
	for id, input := range b.routing.Inputs {
		//if we have already received a value on this input, skip.
		if _, ok := b.state.inputValues[RouteIndex(id)]; ok {
			continue
		}

		// if there is a value set for this input, place value on
		// buffer and set it in map.
		query, ok := input.Value.(*fetch.Query)
		if !ok {
			b.state.inputValues[RouteIndex(id)] = Copy(input.Value)
			continue
		}

		select {
		case m := <-input.C:
			b.state.inputValues[RouteIndex(id)], err = fetch.Run(query, m)
			if err != nil {
				log.Fatal(err)
			}
		case f := <-b.routing.InterruptChan:
			return f
		}
	}
	return nil
}

// run kernel on inputs, produce outputs
func (b *Block) process() Interrupt {
	if b.state.Processed == true {
		return nil
	}

	// if this kernel relies on an external shared state then we need to
	// block until an interrupt connects us to a shared external state.

	if b.sourceType != NONE && b.routing.Source == nil {
		select {
		case f := <-b.routing.InterruptChan:
			return f
		}
	}

	// we should only be able to get here if
	// - we don't need an shared state
	// - we have an external shared state and it has been attached
	if b.sourceType != NONE {
		b.routing.Source.Lock()
	}

	// run the kernel
	interrupt := b.kernel(b.state.inputValues,
		b.state.outputValues,
		b.state.internalValues,
		b.routing.Source,
		b.routing.InterruptChan)

	if interrupt != nil {
		if b.sourceType != NONE {
			b.routing.Source.Unlock()
		}
		return interrupt
	}

	if b.sourceType != NONE {
		b.routing.Source.Unlock()
	}

	b.state.Processed = true

	return nil
}

// broadcast the kernel output to all connections on all outputs.
func (b *Block) broadcast() Interrupt {
	for id, out := range b.routing.Outputs {
		// if the output key is not present in the output map, then we
		// don't deliver any message
		_, ok := b.state.outputValues[RouteIndex(id)]
		if !ok {
			continue
		}

		// if there no connection for this output then wait until there
		// is one. that means we have to wait for an interrupt.
		if len(out.Connections) == 0 {
			select {
			case f := <-b.routing.InterruptChan:
				return f
			}
		}
		for c, _ := range out.Connections {
			// check to see if we have delivered a message to this
			// connection for this block crank. if we have, then
			// skip this delivery.
			m := ManifestPair{id, c}
			if _, ok := b.state.manifest[m]; ok {
				continue
			}

			select {
			case c <- b.state.outputValues[RouteIndex(id)]:
				// set that we have delivered the message.
				b.state.manifest[m] = struct{}{}
			case f := <-b.routing.InterruptChan:
				return f
			}
		}

	}
	return nil
}

// cleanup all block state for this crank of the block
func (b *Block) crank() {
	for k, _ := range b.state.inputValues {
		delete(b.state.inputValues, k)
	}
	for k, _ := range b.state.outputValues {
		delete(b.state.outputValues, k)
	}
	for k, _ := range b.state.manifest {
		delete(b.state.manifest, k)
	}
	b.state.Processed = false
}

package core

import (
	"errors"
	"log"
	"time"
)

// NewBlock creates a new block from a spec
func NewBlock(s Spec) *Block {
	var in []Input
	var out []Output

	for _, v := range s.Inputs {
		in = append(in, Input{
			Name:  v.Name,
			Type:  v.Type,
			Value: nil,
			C:     make(chan Message, 1),
		})
	}

	for _, v := range s.Outputs {
		out = append(out, Output{
			Name:        v.Name,
			Type:        v.Type,
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
		Monitor:    make(chan MonitorMessage, 1),
		lastCrank:  time.Now(),
		done:       make(chan struct{}),
	}
}

func (b *Block) Serve() {
	defer func() {
		b.done <- struct{}{}
	}()
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
			b.routing.Unlock()
			return
		}
		b.routing.Unlock()
	}
}

func (b *Block) exportInput(id RouteIndex) (*Input, error) {
	if int(id) >= len(b.routing.Inputs) || int(id) < 0 {
		return nil, errors.New("index out of range")
	}

	if b.routing.Inputs[id].Value == nil {
		return &b.routing.Inputs[id], nil
	}

	return &Input{
		Value: &InputValue{
			Data: Copy((*b.routing.Inputs[id].Value).Data),
		},
		C:    b.routing.Inputs[id].C,
		Name: b.routing.Inputs[id].Name,
	}, nil

}

// GetInput returns the specified Input
func (b *Block) GetInput(id RouteIndex) (Input, error) {
	b.routing.RLock()
	r, err := b.exportInput(id)
	if err != nil {
		b.routing.RUnlock()
		return Input{}, err
	}
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

// RouteValue sets the route to always be the specified value
func (b *Block) SetInput(id RouteIndex, v *InputValue) error {
	returnVal := make(chan error, 1)
	b.routing.InterruptChan <- func() bool {
		if int(id) < 0 || int(id) >= len(b.routing.Inputs) {
			returnVal <- errors.New("input out of range")
			return true
		}

		// if our receive() has already set the inputValue for the kernel
		// then delete the value out of the input map and use the new one
		if _, ok := b.state.inputValues[id]; ok {
			delete(b.state.inputValues, id)
		}

		b.routing.Inputs[id].Value = v

		returnVal <- nil
		return true
	}
	return <-returnVal
}

// Outputs return a list of manifest pairs for the block
func (b *Block) GetOutputs() []Output {
	b.routing.RLock()
	m := make([]Output, len(b.routing.Outputs), len(b.routing.Outputs))
	for id, out := range b.routing.Outputs {
		m[id] = Output{
			Name:        out.Name,
			Type:        out.Type,
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
			returnVal <- errors.New("invalid source type for this block")
			return true
		}
		b.routing.Source = s
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

func (b *Block) Reset() {
	b.crank()

	// reset block's state as well. currently this only applies to a handful of
	// blocks, like GET and first.
	for k, _ := range b.state.internalValues {
		delete(b.state.internalValues, k)
	}

	// if there are any messages on the input channels, flush them.
	// note: all blocks that are sending to this block MUST BE IN A
	// STOPPED STATE. if any block routines that posess this block's
	// input channel are in a RUNNING state, this flush will not work
	// because it will simply pull another message into the buffer.
	for _, input := range b.routing.Inputs {
		select {
		case <-input.C:
		default:
		}
	}

	return
}

func (b *Block) Stop() {
	b.routing.InterruptChan <- func() bool {
		return false
	}
	<-b.done
	return
}

// wait and listen for all kernel inputs to be filled.
func (b *Block) receive() Interrupt {
	for id, input := range b.routing.Inputs {
		b.Monitor <- MonitorMessage{
			BI_INPUT,
			id,
		}

		//if we have already received a value on this input, skip.
		if _, ok := b.state.inputValues[RouteIndex(id)]; ok {
			continue
		}

		if input.Value != nil {
			b.state.inputValues[RouteIndex(id)] = Copy(input.Value.Data)
			continue
		}

		select {
		case m := <-input.C:
			b.state.inputValues[RouteIndex(id)] = m
		case f := <-b.routing.InterruptChan:
			return f
		}
	}
	return nil
}

// run kernel on inputs, produce outputs
func (b *Block) process() Interrupt {

	b.Monitor <- MonitorMessage{
		BI_KERNEL,
		nil,
	}

	if b.state.Processed == true {
		return nil
	}

	// block until connected to source if necessary

	if b.sourceType != NONE && b.routing.Source == nil {
		select {
		case f := <-b.routing.InterruptChan:
			return f
		}
	}

	// we should only be able to get here if
	// - we don't need an shared state
	// - we have an external shared state and it has been attached

	// if we have a store, lock it
	// we will use the s interface after the kernel as a flag to indicate
	// that this block is attached to a store
	var s interface{}
	if isStore(b.sourceType) {
		s = b.routing.Source
		store, ok := s.(Store)
		if !ok {
			log.Fatal(s)
		}
		store.Lock()
	}

	// run the kernel
	interrupt := b.kernel(b.state.inputValues,
		b.state.outputValues,
		b.state.internalValues,
		b.routing.Source,
		b.routing.InterruptChan)

	// unlock the store if necessary
	if s != nil {
		store, ok := s.(Store)
		if !ok {
			log.Fatal(s)
		}
		store.Unlock()
	}

	// if an interrupt was receieved, return it
	if interrupt != nil {
		return interrupt
	}

	b.state.Processed = true
	return nil
}

// broadcast the kernel output to all connections on all outputs.
func (b *Block) broadcast() Interrupt {
	for id, out := range b.routing.Outputs {
		b.Monitor <- MonitorMessage{
			BI_OUTPUT,
			id,
		}

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

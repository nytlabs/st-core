package core

import (
	"log"

	"github.com/nikhan/go-fetch"
)

// NewBlock creates a new block from a spec
func NewBlock(s Spec) *Block {
	var in []Route
	var out []Output

	for _, v := range s.Inputs {

		q, _ := fetch.Parse(".")
		in = append(in, Route{
			Name: v.Name,
			Path: q,
			C:    make(chan Message),
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
			make(Manifest),
			false,
		},
		routing: BlockRouting{
			Inputs:        in,
			Outputs:       out,
			InterruptChan: make(chan Interrupt),
			Shared: SharedStore{
				Type: s.Shared,
			},
		},
		kernel: s.Kernel,
	}
}

func (b *Block) process() Interrupt {
	if b.state.Processed == true {
		return nil
	}

	if b.routing.Shared.Type != NONE && b.routing.Shared.Store == nil {
		select {
		case f := <-b.routing.InterruptChan:
			return f
		}
	}

	if b.routing.Shared.Type != NONE {
		b.routing.Shared.Store.Lock()
	}

	interrupt := b.kernel(b.state.inputValues,
		b.state.outputValues,
		b.routing.Shared.Store,
		b.routing.InterruptChan)

	if interrupt != nil {
		if b.routing.Shared.Type != NONE {
			b.routing.Shared.Store.Unlock()
		}
		return interrupt
	}

	if b.routing.Shared.Type != NONE {
		b.routing.Shared.Store.Unlock()
	}

	b.state.Processed = true

	return nil
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

// todo: proper getter/setters of route properties
// 	GetRouteValue, GetRoutePath, GetRouteChan

// Input returns the specfied Route
func (b *Block) Input(id RouteID) Route {
	b.routing.RLock()
	defer b.routing.RUnlock()
	return b.routing.Inputs[id]
}

func (b *Block) Store(s Store) {
	b.routing.InterruptChan <- func() bool {
		b.routing.Shared.Store = s
		return true
	}
}

// RouteValue sets the route to always be the specified value
func (b *Block) RouteValue(id RouteID, v Message) {
	b.routing.InterruptChan <- func() bool {
		b.routing.Inputs[id].Value = &v
		return true
	}
}

// RoutePath sets a Route's Path to the supplied go-fetch Query
func (b *Block) RoutePath(id RouteID, p *fetch.Query) {
	b.routing.InterruptChan <- func() bool {
		b.routing.Inputs[id].Path = p
		b.routing.Inputs[id].Value = nil
		return true
	}
}

// Connect connects a Route, specified by ID, to a connection
func (b *Block) Connect(id RouteID, c Connection) {
	b.routing.InterruptChan <- func() bool {
		b.routing.Outputs[id].Connections[c] = struct{}{}
		return true
	}
}

// Disconnect removes a connection from a Route
func (b *Block) Disconnect(id RouteID, c Connection) {
	b.routing.InterruptChan <- func() bool {
		delete(b.routing.Outputs[id].Connections, c)
		return true
	}
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
		if _, ok := b.state.inputValues[RouteID(id)]; ok {
			continue
		}

		// if there is a value set for this input, place value on
		// buffer and set it in map.
		if input.Value != nil {
			b.state.inputValues[RouteID(id)] = *input.Value
			continue
		}

		select {
		case m := <-input.C:
			b.state.inputValues[RouteID(id)], err = fetch.Run(input.Path, m)
			if err != nil {
				log.Fatal(err)
			}
		case f := <-b.routing.InterruptChan:
			return f
		}
	}
	return nil
}

// broadcast the kernel output to all connections on all outputs.
func (b *Block) broadcast() Interrupt {
	for id, out := range b.routing.Outputs {
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
			m := ManifestPair{out.Name, c}
			if _, ok := b.state.manifest[m]; ok {
				continue
			}

			select {
			case c <- b.state.outputValues[RouteID(id)]:
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

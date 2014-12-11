package main

import (
	"github.com/nikhan/go-fetch"
	"log"
)

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
		},
		routing: BlockRouting{
			Inputs:        in,
			Outputs:       out,
			InterruptChan: make(chan InterruptFunc),
		},
		kernel: s.Kernel,
	}
}

func (b *Block) Serve() {
	for {
		var interrupt InterruptFunc

		b.routing.RLock()
		for {
			interrupt = b.receive()
			if interrupt != nil {
				break
			}
			interrupt = b.kernel(b.state.inputValues, b.state.outputValues, b.routing.InterruptChan)
			if interrupt != nil {
				break
			}
			interrupt = b.broadcast()
			if interrupt != nil {
				break
			}
			// we've successfully completed one full loop
			// empty input buffer
			for k, _ := range b.state.inputValues {
				delete(b.state.inputValues, k)
			}
		}
		b.routing.RUnlock()
		b.routing.Lock()
		if ok := interrupt(); !ok {
			return
		}
		b.routing.Unlock()
	}
}

func (b *Block) Input(id RouteID) Route {
	b.routing.RLock()
	defer b.routing.RUnlock()
	return b.routing.Inputs[id]
}

func (b *Block) RouteValue(id RouteID, v Message) {
	b.routing.InterruptChan <- func() bool {
		b.routing.Inputs[id].Value = &v
		return true
	}
}

func (b *Block) RoutePath(id RouteID, p *fetch.Query) {
	b.routing.InterruptChan <- func() bool {
		b.routing.Inputs[id].Path = p
		b.routing.Inputs[id].Value = nil
		return true
	}
}

func (b *Block) Connect(id RouteID, c Connection) {
	b.routing.InterruptChan <- func() bool {
		b.routing.Outputs[id].Connections[c] = struct{}{}
		return true
	}
}

func (b *Block) Disconnect(id RouteID, c Connection) {
	b.routing.InterruptChan <- func() bool {
		delete(b.routing.Outputs[id].Connections, c)
		return true
	}
}

func (b *Block) Stop() {
	b.routing.InterruptChan <- func() bool {
		return false
	}
}

func (b *Block) receive() InterruptFunc {
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

func (b *Block) broadcast() InterruptFunc {
	for id, out := range b.routing.Outputs {
		// if there no connection for this output then wait until there is one
		// that means we have to wait for an interrupt.
		if len(out.Connections) == 0 {
			select {
			case f := <-b.routing.InterruptChan:
				return f
			}
		}
		for c, _ := range out.Connections {
			select {
			case c <- b.state.outputValues[RouteID(id)]:
			case f := <-b.routing.InterruptChan:
				return f
			}
		}

	}
	return nil
}

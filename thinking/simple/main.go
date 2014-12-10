package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/nikhan/go-fetch"
	"github.com/thejerf/suture"
)

type RouteID int
type Message interface{}

type Route struct {
	Name  string
	Path  *fetch.Query
	Value *Message
	C     chan Message
}

type Connection chan Message

type Output struct {
	Name        string
	Connections map[Connection]struct{}
}

type MessageMap map[RouteID]Message

type InterruptFunc func() bool
type KernelFunc func(MessageMap, MessageMap, chan InterruptFunc) InterruptFunc

type Pin struct {
	Name string
}

type Spec struct {
	Inputs  []Pin
	Outputs []Pin
	Kernel  KernelFunc
}

var Library = map[string]Spec{
	"plus": Spec{
		Inputs: []Pin{
			Pin{
				"addend",
			},
			Pin{
				"addend",
			},
		},
		Outputs: []Pin{
			Pin{
				"sum",
			},
		},
		Kernel: func(in MessageMap, out MessageMap, i chan InterruptFunc) InterruptFunc {
			out[0] = in[0].(float64) + in[1].(float64)
			return nil
		},
	},
	"delay": Spec{
		Inputs: []Pin{
			Pin{
				"passthrough",
			},
			Pin{
				"duration",
			},
		},
		Outputs: []Pin{
			Pin{
				"passthrough",
			},
		},
		Kernel: func(in MessageMap, out MessageMap, i chan InterruptFunc) InterruptFunc {
			t, err := time.ParseDuration(in[1].(string))
			if err != nil {
				out[0] = err
				return nil
			}
			timer := time.NewTimer(t)
			select {
			case <-timer.C:
				out[0] = in[0]
				return nil
			case f := <-i:
				return f
			}
			return nil
		},
	},
	"set": Spec{
		Inputs: []Pin{
			Pin{
				"key",
			},
			Pin{
				"value",
			},
		},
		Outputs: []Pin{
			Pin{
				"object",
			},
		},
		Kernel: func(in MessageMap, out MessageMap, i chan InterruptFunc) InterruptFunc {
			out[0] = map[string]interface{}{
				in[0].(string): in[1],
			}
			return nil
		},
	},
	"log": Spec{
		Inputs: []Pin{
			Pin{
				"log",
			},
		},
		Outputs: []Pin{},
		Kernel: func(in MessageMap, out MessageMap, i chan InterruptFunc) InterruptFunc {
			o, err := json.Marshal(in[0])
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(o))
			return nil
		},
	},
}

type Block struct {
	Inputs        []Route
	Outputs       []Output
	Kernel        KernelFunc
	inputValues   MessageMap
	outputValues  MessageMap
	InterruptChan chan InterruptFunc
	Quit          chan struct{}
}

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
		Inputs:        in,
		Outputs:       out,
		Kernel:        s.Kernel,
		inputValues:   make(MessageMap),
		outputValues:  make(MessageMap),
		InterruptChan: make(chan InterruptFunc),
	}
}

func (b *Block) Serve() {
	for {
		var interrupt InterruptFunc
		for {
			interrupt = b.Receive()
			if interrupt != nil {
				//				fmt.Println("recieve interrupt")
				break
			}
			interrupt = b.Kernel(b.inputValues, b.outputValues, b.InterruptChan)
			if interrupt != nil {
				//				fmt.Println("kernel interrupt")
				break
			}
			interrupt = b.Broadcast()
			if interrupt != nil {
				//				fmt.Println("broadcast interupt")
				break
			}
			// we've successfully completed one full loop
			// empty input buffer
			for k, _ := range b.inputValues {
				delete(b.inputValues, k)
			}
		}
		if ok := interrupt(); !ok {
			return
		}
	}
}

func (b *Block) RouteValue(id RouteID, v Message) {
	b.InterruptChan <- func() bool {
		b.Inputs[id].Value = &v
		return true
	}
}

func (b *Block) RoutePath(id RouteID, p *fetch.Query) {
	b.InterruptChan <- func() bool {
		b.Inputs[id].Path = p
		b.Inputs[id].Value = nil
		return true
	}
}

func (b *Block) Connect(id RouteID, c Connection) {
	b.InterruptChan <- func() bool {
		b.Outputs[id].Connections[c] = struct{}{}
		return true
	}
}

func (b *Block) Disconnect(id RouteID, c Connection) {
	b.InterruptChan <- func() bool {
		delete(b.Outputs[id].Connections, c)
		return true
	}
}

func (b *Block) Stop() {
	b.InterruptChan <- func() bool {
		return false
	}
}

func (b *Block) Receive() InterruptFunc {
	var err error
	for id, input := range b.Inputs {
		//if we have already received a value on this input, skip.
		if _, ok := b.inputValues[RouteID(id)]; ok {
			continue
		}

		// if there is a value set for this input, place value on
		// buffer and set it in map.
		if input.Value != nil {
			b.inputValues[RouteID(id)] = *input.Value
			continue
		}

		select {
		case m := <-input.C:
			b.inputValues[RouteID(id)], err = fetch.Run(input.Path, m)
			if err != nil {
				log.Fatal(err)
			}
		case f := <-b.InterruptChan:
			return f
		}
	}
	return nil
}

func (b *Block) Broadcast() InterruptFunc {
	for id, out := range b.Outputs {
		// if there no connection for this output then wait until there is one
		// that means we have to wait for an interrupt.
		if len(out.Connections) == 0 {
			select {
			case f := <-b.InterruptChan:
				return f
			}
		}
		for c, _ := range out.Connections {
			select {
			case c <- b.outputValues[RouteID(id)]:
			case f := <-b.InterruptChan:
				return f
			}
		}

	}
	return nil
}

func main() {
	supervisor := suture.NewSimple("st-core")
	supervisor.ServeBackground()

	b := NewBlock(Library["plus"])
	d := NewBlock(Library["delay"])
	l := NewBlock(Library["log"])
	o := NewBlock(Library["set"])
	l2 := NewBlock(Library["log"])

	a := supervisor.Add(b)
	_ = supervisor.Add(d)
	_ = supervisor.Add(l)
	_ = supervisor.Add(o)
	_ = supervisor.Add(l2)

	b.Connect(0, d.Inputs[0].C)
	d.Connect(0, o.Inputs[1].C)
	o.Connect(0, l.Inputs[0].C)

	path, _ := fetch.Parse(".test")
	fmt.Println(path)

	b.RouteValue(RouteID(0), 1.1)
	d.RouteValue(RouteID(1), "10ms")
	o.RouteValue(RouteID(0), "test")

	go func() {
		for {
			b.RouteValue(RouteID(1), rand.Float64()*10.0)
			time.Sleep(time.Duration(rand.Intn(2)+1) * time.Millisecond)
		}

	}()

	go func() {
		for {
			_ = b.Inputs[0].Path
			time.Sleep(1 * time.Millisecond)
		}
	}()

	time.Sleep(500 * time.Millisecond)

	fmt.Println("Disconnected!")

	d.Disconnect(0, o.Inputs[1].C)

	l2.RoutePath(RouteID(0), path)
	o.Connect(RouteID(0), l2.Inputs[0].C)

	time.Sleep(300 * time.Millisecond)

	d.Connect(0, o.Inputs[1].C)
	fmt.Println("Connected")

	time.Sleep(100 * time.Millisecond)

	supervisor.Remove(a)

	supervisor.Stop()
}

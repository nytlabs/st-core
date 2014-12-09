package main

import (
	"log"
	"time"

	"github.com/nikhan/go-fetch"
	"github.com/thejerf/suture"
)

const (
	QUIT = iota
	ADD_CONNECTION
	DEL_CONNECTION
	SET_PATH
	SET_VALUE
)

type RouteID int64
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

type RxKernel func(MessageMap, MessageMap, chan struct{}) bool

type Block struct {
	Inputs  []Route
	Outputs []Output
	Kernel  RxKernel
	Control chan interface{}
	Quit    chan struct{}
}

func NewBlock() *Block {
	return &Block{
		Inputs:  []Route{},
		Outputs: []Output{},
		Control: make(chan interface{}),
		Quit:    make(chan struct{}),
	}
}

func (b *Block) Serve() {
	inputs := make(MessageMap)
	outputs := make(MessageMap)
	for {
		if ok := b.Receive(inputs); !ok {
			break
		}

		if ok = b.Kernel(inputs, outputs, b.QuitChan); !ok {
			break
		}

		if ok = b.Broadcast(outputs); !ok {
			break
		}

		// we've successfully completed one full loop
		// empty input buffer
		for k, _ := range inputs {
			delete(inputs, k)
		}
	}
}

func (b *Block) SetValue(id RouteId, value Message) {
	b.Break <- struct{}{}
}

func (b *Block) SetPath(id RouteId, value *fetch.Query) {
	b.Break <- struct{}{}
}

func (b *Block) Connect(id RouteId, c Connection) {
	b.Break <- struct{}{}
}

func (b *Block) Disconnect(id RouteId, c Connection) {
	b.Break <- struct{}{}
}

func (b *Block) Stop() {
	b.Break <- struct{}{}
}

func (b *Block) Receive(values MessageMap) bool {
	var err Error
	for id, input := range b.Inputs {
		//if we have already received a value on this input, skip.
		if _, ok := values[id]; ok {
			continue
		}

		// if there is a value set for this input, place value on
		// buffer and set it in map.
		if input.Value != nil {
			values[id] = *input.Value
			continue
		}

		select {
		case m := <-input.C:
			v, err = fetch.Run(input.Path, m)
			if err != nil {
				log.Fatal(err)
			}
			values[id] = v
		case <-b.QuitChan:
			return nil, false
		}
	}
	return values, true
}

func (b *Block) Broadcast(values MessageMap) bool {
	for _, _ := range values {
		/*o := b.GetOutput(k)
		for c, _ := range o.GetConnections() {
			select {
			case c <- v:
			case <-b.QuitChan:
				return false
			}
		}*/
	}
	return true
}

func main() {
	supervisor := suture.NewSimple("st-core")
	supervisor.ServeBackground()

	a := supervisor.Add(NewBlock())

	time.Sleep(1 * time.Second)

	supervisor.Remove(a)

	supervisor.Stop()
}

package main

import (
	"log"
	"time"

	"github.com/nikhan/go-fetch"
	"github.com/thejerf/suture"
)

type RouteID int64
type Message interface{}

type Route struct {
	Name  string
	Path  *fetch.Query
	Value *Message
	C     chan Message
}

func (r *Route) SetValue(m Message) {
	r.Value = m
}

type MessageMap map[RouteID]Message

type RxKernel func(MessageMap, MessageMap, chan struct{}) bool

type Block struct {
	Inputs  []Route
	Outputs []Route
	Kernel  RxKernel
	Quit    chan struct{}
}

func NewBlock() *Block {
	return &Block{
		Inputs:  []Route{},
		Outputs: []Route{},
		Quit:    make(chan struct{}),
	}
}

func (b *Block) Serve() {
	inValues := make(MessageMap)
	outValues := make(MessageMap)

	for {
		if ok := b.Receive(inValues); !ok {
			break
		}

		if ok = b.Kernel(inValues, outValues, b.QuitChan); !ok {
			break
		}

		if ok = b.Broadcast(outValues); !ok {
			break
		}
	}
}

func (b *Block) Receive(values MessageMap) bool {
	var err Error
	for id, input := range b.Inputs {
		if input.Value != nil {
			values[id] = *input.Value
			continue
		}

		select {
		case m := <-input.C:
			values[id], err = fetch.Run(input.Path, m)
			if err != nil {
				log.Fatal(err)
			}
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

func (b *Block) Stop() {
	b.Quit <- struct{}{}
}

func main() {
	supervisor := suture.NewSimple("st-core")
	supervisor.ServeBackground()

	a := supervisor.Add(NewBlock())

	time.Sleep(1 * time.Second)

	supervisor.Remove(a)

	supervisor.Stop()
}

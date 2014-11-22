package core

import (
	"log"

	"github.com/nikhan/go-fetch"
)

type F struct {
	*Block
}

func NewF(name string) F {
	b := NewBlock(name)
	b.AddInput("in")
	b.AddOutput("out")
	return F{b}
}

func (b F) Kernel() {
	log.Println("F")
}

func (b F) Serve() {
	in := b.GetInput("in")
	var msg Message
	var err error

	for {

		select {
		case m := <-in.Connection:
			msg, err = fetch.Run(in.Path, m)
			if err != nil {
				log.Fatal(err)
			}
		case <-b.QuitChan:
			return
		}

		b.Kernel()

		if ok := b.Broadcast(msg, "out"); !ok {
			return
		}

	}
}

func (b F) String() string {
	return "F"
}

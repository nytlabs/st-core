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

	b.Kernel = func(msgs ...Message) (map[string]Message, error) {
		log.Println(msgs)
		return map[string]Message{
			"out": msgs[0],
		}, nil
	}
	return F{b}
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

		b.Kernel(msg)

	}
}

func (b F) String() string {
	return "F"
}

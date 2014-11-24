package core

import (
	"log"

	"github.com/nikhan/go-fetch"
)

type G struct {
	*Block
}

func NewG(name string) G {
	b := NewBlock(name)
	b.AddInput("in")
	b.AddOutput("out")
	b.Kernel = func(msgs ...Message) (map[string]Message, error) {
		log.Println(msgs)
		return map[string]Message{
			"out": msgs[0],
		}, nil
	}
	return G{b}
}

func (b G) Kernel() {
	log.Println("G")
}

func (b G) Serve() {
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

func (b G) String() string {
	return "G"
}

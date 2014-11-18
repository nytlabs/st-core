package core

import (
	"errors"
	"log"

	"github.com/nikhan/go-fetch"
)

// The Set block maps an inbound message onto an outbound message using the supplied rule.
// Them Set block has two inputs: "in" and "mapping", and one outupt: "out".
type Set struct {
	*Block
}

func NewSet(name string) Set {
	b := NewBlock(name)
	b.AddInput("in")
	b.AddInput("key")
	b.AddOutput("out")
	return Set{b}
}

func (b Set) Serve() {

	in := b.GetInput("in")
	key := b.GetInput("key")

	var k string
	var msg interface{}
	var err error
	var ok bool

	for {
		select {
		// get the message
		case m := <-in.Connection:
			msg, err = fetch.Run(in.Path, m)
			if err != nil {
				log.Fatal(err)
			}
		case <-b.QuitChan:
			return
		}

		// get the key
		select {
		case m := <-key.Connection:
			kI, err := fetch.Run(key.Path, m)
			if err != nil {
				log.Fatal(err)
			}
			k, ok = kI.(string)
			if !ok {
				log.Fatal(errors.New("supplied key must be a string"))
			}
		case <-b.QuitChan:
			return
		}

		out := map[string]interface{}{
			k: msg,
		}

		if ok := b.Broadcast(out, "out"); !ok {
			return
		}
	}
}

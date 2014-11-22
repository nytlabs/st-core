package core

import (
	"fmt"
	"log"

	"github.com/nikhan/go-fetch"
)

// Pusher applies constant pressure to its outbound Route. If a value is supplied that value is pushed, otherwise a nil messge is pushed. A pusher block is special in that it will always provide pressure despite lack of pressure on its input.
type Pusher struct {
	*Block
}

func NewPusher(name string) Pusher {
	b := NewBlock(name)
	b.AddInput("value")
	b.AddOutput("out")
	return Pusher{
		b,
	}
}

func (b Pusher) Serve() {

	value := b.GetInput("value")

	var msg Message
	var err error

	for {
		// get the value to be pushed
		select {
		case m := <-value.Connection:
			msg, err = fetch.Run(value.Path, m)
			if err != nil {
				log.Fatal(err)
			}
		case <-b.QuitChan:
			return
		default:
			// default is the nil Message
		}
		if ok := b.Broadcast(msg, "out"); !ok {
			return
		}

	}
}

func (b Pusher) String() string {
	return fmt.Sprintf("Pusher: %s", b.Name)
}

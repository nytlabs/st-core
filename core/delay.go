package core

import (
	"log"
	"time"

	"github.com/nikhan/go-fetch"
)

type Delay struct {
	*Block
}

func NewDelay() Delay {
	b := NewBlock("delay")
	b.AddInput("in")
	b.AddOutput("out")
	return Delay{b}
}

func (b Delay) Serve() {
	in := b.GetInput("in")
	var msg Message
	var err error

	for {
		t := time.NewTimer(1 * time.Second)
		select {
		case <-t.C:
		case <-b.QuitChan:
			return
		}

		select {
		case m := <-in.Connection:
			msg, err = fetch.Run(in.Path, m)
			if err != nil {
				log.Fatal(err)
			}
		case <-b.QuitChan:
			return
		}

		for c, _ := range b.Connections("out") {
			select {
			case c <- msg:
			case <-b.QuitChan:
				return
			}
		}
	}
}

func (b Delay) String() string {
	return "Delay"
}

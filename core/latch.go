package core

import (
	"errors"
	"log"

	"github.com/nikhan/go-fetch"
)

type Latch struct {
	*Block
}

func NewLatch(name string) Latch {
	b := NewBlock(name)
	b.AddInput("in")
	b.AddInput("ctrl")
	b.AddOutput("out")
	return Latch{b}
}

func (b Latch) Serve() {
	in := b.GetInput("in")
	ctrl := b.GetInput("ctrl")

	var err error
	var msg, ctrlSignal Message

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

		select {
		case m := <-ctrl.Connection:
			ctrlSignal, err = fetch.Run(ctrl.Path, m)
			if err != nil {
				log.Fatal(err)
			}
		case <-b.QuitChan:
			return
		}

		switch ctrlSignal := ctrlSignal.(type) {
		case bool:
			if ctrlSignal {
				if ok := b.Broadcast(msg, "out"); !ok {
					return
				}
			}
		case error:
			log.Fatal(ctrlSignal)
		default:
			log.Fatal(errors.New("unrecognised control signal in latch " + b.Name))
		}

	}
}

package core

import (
	//    "fmt"
	"log"

	"github.com/nikhan/go-fetch"
)

type Greater struct {
	*Block
}

func NewGreater(name string) Greater {
	b := NewBlock(name)
	b.AddInput("addend 1")
	b.AddInput("addend 2")
	b.AddOutput("out")
	return Greater{b}
}

func (b Greater) Serve() {
	in1 := b.GetInput("addend 1")
	in2 := b.GetInput("addend 2")

	var add, bdd interface{}
	var err error

	for {

		select {
		case m := <-in1.Connection:
			add, err = fetch.Run(in1.Path, m)
			if err != nil {
				log.Fatal(err)
			}
		case <-b.QuitChan:
			return
		}

		v1, ok := add.(float64)
		if !ok {
			log.Fatal("coudln't assert first input to float")
		}

		select {
		case m := <-in2.Connection:
			bdd, err = fetch.Run(in2.Path, m)
			if err != nil {
				log.Fatal(err)
			}
		case <-b.QuitChan:
			return
		}

		v2, ok := bdd.(float64)
		if !ok {
			log.Fatal("coudln't assert second input to float")
		}

		result := false
		if v1 > v2 {
			result = true
		}

		if ok := b.Broadcast(result, "out"); !ok {
			return
		}
	}
}

package core

import (
	"fmt"
	"log"

	"github.com/nikhan/go-fetch"
)

type Plus struct {
	*Block
}

func NewPlus() Plus {
	b := NewBlock("plus")
	b.AddInput("addend 1")
	b.AddInput("addend 2")
	b.AddOutput("out")
	return Plus{b}
}

func (b Plus) Serve() {
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

		select {
		case m := <-in2.Connection:
			bdd, err = fetch.Run(in2.Path, m)
			if err != nil {
				log.Fatal(err)
			}
		case <-b.QuitChan:
			return
		}

		result := add.(int) + bdd.(int)
		fmt.Println(add, "+", bdd, "=", result)
		for c, _ := range b.Connections("out") {
			select {
			case c <- result:
			case <-b.QuitChan:
				return
			}
		}
	}
}

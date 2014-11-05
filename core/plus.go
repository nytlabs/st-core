package core

import "fmt"

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
	in1 := b.GetInput("addend 1").Connection
	in2 := b.GetInput("addend 2").Connection
	for {
		bdd := <-in2
		add := <-in1
		result := add.(int) + bdd.(int)
		fmt.Println(add, "+", bdd)
		for c, _ := range b.Connections("out") {
			select {
			case c <- result:
			case <-b.QuitChan:
				return
			}
		}
	}
}

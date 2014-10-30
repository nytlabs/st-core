package core

import "fmt"

type Plus struct {
	*Block
}

func NewPlus() Plus {
	b := NewBlock()
	b.AddInput("addend 1")
	b.AddInput("addend 2")
	b.AddOutput("out")
	return Plus{b}
}

func (b Plus) Serve() {
	in1 := b.GetInput("addend 1")
	in2 := b.GetInput("addend 2")
	for {
		bdd := <-in2
		add := <-in1
		c := add.(int) + bdd.(int)
		fmt.Println(add, "+", bdd)
		b.Broadcast("out", c)
	}
}

func (b Plus) Stop() {
	for c, _ := range b.Outputs["out"].Connections {
		close(c)
	}
}

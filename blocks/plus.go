package blocks

import (
	"fmt"

	"github.com/nytlabs/st-core/core"
)

type Plus struct {
	*core.Block
}

func NewPlus() Plus {
	b := core.NewBlock()
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

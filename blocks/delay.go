package blocks

import (
	"time"

	"github.com/nytlabs/st-core/core"
)

type Delay struct {
	*core.Block
}

func NewDelay() Delay {
	b := core.NewBlock()
	b.AddInput("in")
	b.AddOutput("out")
	return Delay{b}
}

func (b Delay) Serve() {
	for {
		m := <-b.GetInput("in")
		time.Sleep(10 * time.Millisecond)
		b.Broadcast("out", m)
	}
}

func (b Delay) Stop() {
	for c, _ := range b.Outputs["out"].Connections {
		close(c)
	}
}

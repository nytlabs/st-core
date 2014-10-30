package blocks

import "github.com/nytlabs/st-core/core"

type Pusher struct {
	*core.Block
}

func NewPusher() Pusher {
	b := core.NewBlock()
	b.AddOutput("out")
	return Pusher{
		b,
	}
}

func (b Pusher) Serve() {
	i := 0
	for {
		i++
		b.Broadcast("out", i)
	}
}

func (b Pusher) Stop() {
	for c, _ := range b.Outputs["out"].Connections {
		close(c)
	}
}

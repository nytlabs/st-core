package blocks

import (
	"fmt"

	"github.com/nytlabs/st-core/core"
)

type Log struct {
	*core.Block
	name string
}

func NewLog(name string) Log {
	b := core.NewBlock()
	b.AddInput("in")
	return Log{
		Block: b,
		name:  name,
	}
}

func (b Log) Serve() {
	for {
		m := <-b.GetInput("in")
		fmt.Println(b.name, ": ", m)
	}
}
func (b Log) Stop() {
	for c, _ := range b.Outputs["out"].Connections {
		close(c)
	}
}

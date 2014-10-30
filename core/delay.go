package core

import "time"

type Delay struct {
	Block
}

func NewDelay() Delay {
	b := NewBlock()
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

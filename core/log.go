package core

import "fmt"

type Log struct {
	*Block
	name string
}

func NewLog(name string) Log {
	b := NewBlock()
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

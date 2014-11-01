package core

import "fmt"

type Log struct {
	*Block
	name string
}

func NewLog(name string) Log {
	b := NewBlock(name)
	b.AddInput("in")
	return Log{
		Block: b,
	}
}

func (b Log) Serve() {
	for {
		select {
		case m := <-b.GetInput("in"):
			fmt.Println(b.Name, ": ", m)
		case <-b.QuitChan:
			return
		}
	}
}

func (b Log) String() string {
	return "Log"
}

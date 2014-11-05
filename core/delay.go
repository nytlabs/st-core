package core

import "time"

type Delay struct {
	*Block
}

func NewDelay() Delay {
	b := NewBlock("delay")
	b.AddInput("in")
	b.AddOutput("out")
	return Delay{b}
}

func (b Delay) Serve() {
	for {
		m := <-b.GetInput("in").Connection
		time.Sleep(1 * time.Second)

		for c, _ := range b.Connections("out") {
			select {
			case c <- m:
			case <-b.QuitChan:
				return
			}
		}
	}
}

func (b Delay) String() string {
	return "Delay"
}

package core

import (
	"fmt"
	"time"
)

// Pusher applies constant pressure to its outbound Route
type Pusher struct {
	*Block
}

func NewPusher(name string) Pusher {
	b := NewBlock(name)
	b.AddOutput("out")
	return Pusher{
		b,
	}
}

func (b Pusher) Serve() {
	i := 0
	for {
		i++
		for c := range b.Connections("out") {
			select {
			case c <- i:
			case <-b.QuitChan:
				time.Sleep(120 * time.Second)
				return
			}
		}
	}
}

func (b Pusher) String() string {
	return fmt.Sprintf("Pusher: %s", b.Name)
}

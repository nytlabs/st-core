package core

type Pusher struct {
	*Block
}

func NewPusher() Pusher {
	b := NewBlock()
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

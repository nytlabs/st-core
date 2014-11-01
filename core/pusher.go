package core

type Pusher struct {
	*Block
}

func NewPusher() Pusher {
	b := NewBlock("pusher")
	b.AddOutput("out")
	return Pusher{
		b,
	}
}

func (b Pusher) Serve() {
	i := 0
	for {
		i++
		// broadcast
		b.Broadcast("out", i)
	}
}

func (b Pusher) String() string {
	return "Pusher"
}

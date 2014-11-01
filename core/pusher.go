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
		for c, _ := range b.Connections("out") {
			select {
			case c <- i:
			case <-b.QuitChan:
				return
			}
		}
	}
}

func (b Pusher) String() string {
	return "Pusher"
}

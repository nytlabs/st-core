package core

type Plus struct {
	*Block
}

func NewPlus(name string) Plus {
	b := NewBlock("plus")
	b.AddInput("addend 1")
	b.AddInput("addend 2")
	b.AddOutput("out")
	return Plus{b}
}

func (b Plus) Serve() {
	in1 := b.GetInput("addend 1")
	in2 := b.GetInput("addend 2")

	for {
		if ok := b.Recieve(); !ok {
			return
		}

		result := in1.Value.(int) + in2.Value.(int)

		if ok := b.Broadcast(result, "out"); !ok {
			return
		}
	}
}

package core

import (
	"log"
	"sync"

	"github.com/nikhan/go-fetch"
)

// A Block is the basic processing unit in streamtools. It has inbound and outbound routes.
type Block struct {
	Name     string // for logging
	Inputs   map[string]*Input
	Outputs  map[string]*Output
	QuitChan chan bool
	sync.Mutex
	Kernel func(...Message) (map[string]Message, error) // route -> message to be sent
}

// NewBlock returns a block with no inputs and no outputs.
func NewBlock(name string) *Block {
	return &Block{
		Name:     name,
		Inputs:   make(map[string]*Input),
		Outputs:  make(map[string]*Output),
		QuitChan: make(chan bool),
	}
}

// Add a named input to the block
func (b *Block) AddInput(id string) bool {
	b.Lock()
	defer b.Unlock()
	_, ok := b.Inputs[id]
	if ok {
		return false
	}
	b.Inputs[id] = NewInput()
	return true
}

// Remove a named input to the block
func (b *Block) RemoveInput(id string) bool {
	b.Lock()
	defer b.Unlock()
	_, ok := b.Inputs[id]
	if !ok {
		return false
	}
	delete(b.Inputs, id)
	return true
}

// GetInput returns the input Connection
func (b *Block) GetInput(id string) *Input {
	b.Lock()
	input, ok := b.Inputs[id]
	b.Unlock()
	if !ok {
		return nil
	}
	return input
}

// AddOutput registers a new output for the block
func (b *Block) AddOutput(id string) bool {
	b.Lock()
	defer b.Unlock()
	_, ok := b.Outputs[id]
	if ok {
		return false
	}
	b.Outputs[id] = NewOutput()
	return true
}

// RemoveOutput deletes the output from the block
func (b *Block) RemoveOutput(id string) bool {
	b.Lock()
	defer b.Unlock()
	_, ok := b.Outputs[id]
	if !ok {
		return false
	}
	delete(b.Outputs, id)
	return true
}

// Stop is called when removing a block from the streamtools pattern. This is the default, and can be overwritten.
func (b *Block) Stop() {
	b.QuitChan <- true
}

// Broadcast is called when sending a message to an Output. If Broadcast returns false your block must immediately return.
func (b Block) Broadcast(m Message, id string) bool {
	for c, _ := range b.Connections(id) {
		select {
		case c <- m:
		case <-b.QuitChan:
			return false
		}
	}
	return true
}

func (b Block) getName() string {
	return b.Name
}

func (b Block) Merge(β Block) *Block {
	out := NewBlock(b.getName() + "_" + β.getName())
	for id, input := range b.Inputs {
		out.Inputs[id] = input
	}
	for id, output := range β.Outputs {
		out.Outputs[id] = output
	}
	out.Kernel = func(msgs ...Message) (map[string]Message, error) {
		outMsg, err := b.Kernel(msgs)
		if err != nil {
			return nil, err
		}
		inMsg := map[string]Message{
			"in": outMsg["out"],
		}
		return β.Kernel(inMsg)
	}
	return out
}

func (b Block) Recieve() bool {
	var err error
	for _, in := range b.Inputs {
		select {
		case m := <-in.Connection:
			in.Value, err = fetch.Run(in.Path, m)
			if err != nil {
				log.Fatal(err)
			}
		case <-b.QuitChan:
			return false
		}
	}
	return true
}

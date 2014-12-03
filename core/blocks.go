package core

import (
	"log"
	"sync"

	"github.com/nikhan/go-fetch"
)

type KernelFunc func(chan bool, map[string]Message) (map[string]Message, bool)

type Spec struct {
	Name    string
	Inputs  []string
	Outputs []string
	Kernel  KernelFunc
}

// A Block is the basic processing unit in streamtools. It has inbound and outbound routes.
type Block struct {
	Name     string // for logging
	Inputs   map[string]*Input
	Outputs  map[string]*Output
	QuitChan chan bool
	Kernel   KernelFunc
	sync.Mutex
}

// NewBlock returns a block with no inputs and no outputs.
func NewBlock(s Spec) *Block {

	nb := &Block{
		Name:     s.Name,
		Inputs:   make(map[string]*Input),
		Outputs:  make(map[string]*Output),
		QuitChan: make(chan bool),
	}

	for _, v := range s.Inputs {
		nb.AddInput(v)
	}

	for _, v := range s.Outputs {
		nb.AddOutput(v)
	}

	nb.Kernel = s.Kernel

	return nb
}

func (b *Block) Serve() {
	var values map[string]Message
	var output map[string]Message
	var ok bool
	for {
		if values, ok = b.Receive(); !ok {
			return
		}

		if output, ok = b.Kernel(b.QuitChan, values); !ok {
			return
		}

		if ok = b.Broadcast(output); !ok {
			return
		}
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

// GetOutput returns the specified output
func (b *Block) GetOutput(id string) *Output {
	b.Lock()
	output, ok := b.Outputs[id]
	b.Unlock()
	if !ok {
		return nil
	}
	return output
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
func (b Block) Broadcast(outputs map[string]Message) bool {
	for k, v := range outputs {
		o := b.GetOutput(k)
		for c, _ := range o.GetConnections() {
			select {
			case c <- v:
			case <-b.QuitChan:
				return false
			}
		}
	}
	return true
}

func (b Block) getName() string {
	return b.Name
}

func (b Block) Merge(β Block) *Block {

	out := &Block{
		Name: b.getName() + "_" + β.getName(),
	}

	for id, input := range b.Inputs {
		out.Inputs[id] = input
	}
	for id, output := range β.Outputs {
		out.Outputs[id] = output
	}

	out.Kernel = func(quitChan chan bool, msgs map[string]Message) (map[string]Message, bool) {
		outMsg, ok := b.Kernel(quitChan, msgs)
		if !ok {
			return nil, false
		}
		inMsg := map[string]Message{
			"in": outMsg["out"],
		}
		return β.Kernel(quitChan, inMsg)
	}
	return out
}

func (b Block) Receive() (map[string]Message, bool) {
	var err error
	values := make(map[string]Message)
	for name, in := range b.Inputs {
		select {
		case m := <-in.Connection:
			values[name], err = fetch.Run(in.Path, m)
			if err != nil {
				log.Fatal(err)
			}
		case <-b.QuitChan:
			return nil, false
		}
	}
	return values, true
}

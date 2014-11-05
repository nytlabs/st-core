package core

import (
	"sync"

	"github.com/nikhan/go-fetch"
)

// A Message flows through a connection
type Message interface{}

// A Connection passes messages from block to block
type Connection chan Message

// An Output is a collection of Connections
type Output struct {
	sync.Mutex
	Connections map[Connection]bool
}

// An Input owns a single connection that can be shared by multiple Outputs
type Input struct {
	Path       *fetch.Query
	Value      string
	Connection Connection
}

func NewInput() *Input {
	return &Input{
		Connection: make(Connection),
	}
}

// Constructs a new Route with no connections
func NewOutput() *Output {
	return &Output{
		Connections: make(map[Connection]bool),
	}
}

// Add a Connection to an Output
func (r *Output) Add(c Connection) bool {
	_, ok := r.Connections[c]
	if ok {
		return false
	}
	r.Connections[c] = true
	return true
}

// Remove a Connection from a Route
func (r *Output) Remove(c Connection) bool {
	_, ok := r.Connections[c]
	if !ok {
		return false
	}
	delete(r.Connections, c)
	return true
}

// A Block is the basic processing unit in streamtools. It has inbound and outbound routes.
type Block struct {
	Name     string // for logging
	Inputs   map[string]*Input
	Outputs  map[string]*Output
	QuitChan chan bool
	sync.Mutex
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

// Set an input route's Path
func (b *Block) SetPath(id, path string) error {
	query, err := fetch.Parse(path)
	if err != nil {
		return err
	}
	b.Lock()
	b.Inputs[id].Path = query
	b.Unlock()
	return nil
}

// Set an input route's Value
func (b *Block) SetValue(id, value string) error {
	b.Lock()
	b.Inputs[id].Value = value
	b.Unlock()
	return nil
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

// AddOutput registers a new output Route for the block
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

// RemoveOutput deletes the output route from the block
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

// GetConnections returns all the connections associated with the specified output route
func (b *Block) Connections(id string) map[Connection]bool {
	// get route
	b.Lock()
	route := b.Outputs[id]
	b.Unlock()
	// get connections
	route.Lock()
	connections := route.Connections
	route.Unlock()
	return connections
}

// Connect an output Route from this block to a Route elsewhere in streamtools
func (b *Block) Connect(id string, in *Input) bool {
	b.Lock()
	ok := b.Outputs[id].Add(in.Connection)
	b.Unlock()
	return ok
}

// Discconnect an output Route of this block from a previously connected Route
func (b *Block) Disconnect(id string, r Connection) bool {
	b.Lock()
	ok := b.Outputs[id].Remove(r)
	b.Unlock()
	return ok
}

// Stop is called when removing a block from the streamtools pattern.
func (b *Block) Stop() {
	b.QuitChan <- true
}

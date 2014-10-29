package main

import "sync"

// A Connection passes messages from block to block
type Connection chan interface{}

// A Route is a collection of Connections
type Route struct {
	sync.Mutex
	Connections map[Connection]bool
}

// Constructs a new Route with no connections
func NewRoute() *Route {
	return &Route{
		Connections: make(map[Connection]bool),
	}
}

// Add a Connection to a Route
func (r *Route) Add(c Connection) bool {
	_, ok := r.Connections[c]
	if ok {
		return false
	}
	r.Connections[c] = true
	return true
}

// Remove a route from a Route
func (r *Route) Remove(c Connection) bool {
	_, ok := r.Connections[c]
	if !ok {
		return false
	}
	delete(r.Connections, c)
	return true
}

// Send message `m` to all routes controlled by this Route
func (r *Route) Broadcast(m interface{}) {
	r.Lock()
	routes := r.Connections
	r.Unlock()
	for r, _ := range routes {
		r <- m
	}
}

// A Block is the basic processing unit in streamtools. It has inbound and outbound routes.
type Block struct {
	Inputs  map[string]Connection
	Outputs map[string]*Route
	sync.Mutex
}

// NewBlock returns a block with no inputs and no outputs.
func NewBlock() *Block {
	return &Block{
		Inputs:  make(map[string]Connection),
		Outputs: make(map[string]*Route),
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
	b.Inputs[id] = make(Connection)
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

// GetInput returns the input Route
func (b *Block) GetInput(id string) Connection {
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
	b.Outputs[id] = NewRoute()
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

// Connect an output Route from this block to a Route elsewhere in streamtools
func (b *Block) Connect(id string, r Connection) bool {
	b.Lock()
	ok := b.Outputs[id].Add(r)
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

// Broadcast a message from this block to its output Routes
func (b *Block) Broadcast(id string, m interface{}) {
	b.Lock()
	route := b.Outputs[id]
	b.Unlock()
	route.Broadcast(m)
}

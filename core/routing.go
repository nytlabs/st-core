package core

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

// Remove a Connection from a Route
func (r *Route) Remove(c Connection) bool {
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
	Inputs   map[string]Connection
	Outputs  map[string]*Route
	QuitChan chan bool
	sync.Mutex
}

// NewBlock returns a block with no inputs and no outputs.
func NewBlock(name string) *Block {
	return &Block{
		Name:     name,
		Inputs:   make(map[string]Connection),
		Outputs:  make(map[string]*Route),
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

// GetInput returns the input Connection
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

// GetConnections returns all the connections associated with the specified output route
func (b *Block) Connections(id string) map[Connection]bool {
	// get route
	b.Lock()
	route := b.Outputs["out"]
	b.Unlock()
	// get connections
	route.Lock()
	connections := route.Connections
	route.Unlock()
	return connections
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

func (b *Block) Broadcast(id string, m interface{}) error {
	for c, _ := range b.Connections(id) {
		select {
		case c <- m:
		case <-b.QuitChan:
			return nil
		}
	}
	return nil
}

func (b *Block) Stop() {
	b.QuitChan <- true
}

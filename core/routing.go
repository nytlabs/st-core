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
	sync.Mutex
	Path       *fetch.Query // used to extract information from the inbound message
	Value      Message
	Connection Connection // inbound messages arrive on this Connection
	quitChan   chan bool  // used to interrupt the input's value pusher
}

// NewInput creates an input with its single Connection
func NewInput() *Input {
	q, _ := fetch.Parse(".")
	return &Input{
		Path:       q,
		Connection: make(Connection),
		quitChan:   make(chan bool),
	}
}

// Constructs a new Output ready to be connected
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

// Remove a Connection from an Output
func (r *Output) Remove(c Connection) bool {
	_, ok := r.Connections[c]
	if !ok {
		return false
	}
	delete(r.Connections, c)
	return true
}

// Set an input's Path
func (r *Input) SetPath(path string) error {
	query, err := fetch.Parse(path)
	if err != nil {
		return err
	}
	r.Lock()
	r.Path = query
	r.Unlock()
	return nil
}

// call this whenever you want to set a value, or make a new connection
func stopValuePusher(in *Input) {
	select {
	case in.quitChan <- true:
	default:
		// wasn't running (is there a race here?)
	}
}

// Set an input's Value
func (i *Input) SetValue(value Message) error {
	// we store the marshalled value in the Input so we can access it later
	i.Lock()
	i.Value = value
	i.Unlock()

	// then, to set an input to a particular value, we just push
	// that value to that input, as though we had a little pusher block.

	// first kill any existing value pusher
	stopValuePusher(i)

	// then set the pusher going
	go func() {
		for {
			select {
			case i.Connection <- value:
			case <-i.quitChan:
				return
			}
		}
	}()
	return nil
}

// GetConnections returns all the connections associated with the specified output
func (o *Output) GetConnections() map[Connection]bool {
	// get connections
	o.Lock()
	connections := o.Connections
	o.Unlock()
	return connections
}

// Connect an Output from this block to an Input elsewhere in streamtools
func (o *Output) Connect(in *Input) bool {
	stopValuePusher(in)
	o.Lock()
	ok := o.Add(in.Connection)
	o.Unlock()
	return ok
}

// Discconnect an Output of this block from a previously connected Input
func (o *Output) Disconnect(c Connection) bool {
	o.Lock()
	ok := o.Remove(c)
	o.Unlock()
	return ok
}

package core

import (
	"sync"
	"log"
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
	Connection Connection   // inbound messages arrive on this Connection
	quitChan   chan bool    // used to interrupt the input's value pusher
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

// Set an input's Path
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

// GetConnections returns all the connections associated with the specified output
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

// Connect an Output from this block to an Input elsewhere in streamtools
func (b *Block) Connect(id string, in *Input) bool {
	stopValuePusher(in)
	b.Lock()
	ok := b.Outputs[id].Add(in.Connection)
	b.Unlock()
	return ok
}

// Discconnect an Output of this block from a previously connected Input
func (b *Block) Disconnect(id string, r Connection) bool {
	b.Lock()
	ok := b.Outputs[id].Remove(r)
	b.Unlock()
	return ok
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

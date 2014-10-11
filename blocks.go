package main

import "github.com/thejerf/suture"

// Msg defines the individual messages that flow through streamtools
type Msg interface{}

// Route controls inbound or outbound messages from a block
type Route struct {
	C    chan Msg
	Kind string
}

// BlockInterface defines the basic capabilities of a streamtools block.
type BlockInterface interface {
	suture.Service
}

// Block defines the basic properties of a streamtools block
type Block struct {
	Name   string           // friendly name used in the API and UI
	Desc   string           // short description of the block for the UI
	Routes map[string]Route // container for this block's routes
	Quit   chan bool        // channel to indicate this block should quit
}

// NewBlock is the basic constructor for a Block.
func NewBlock(name, desc string) *Block {
	return &Block{
		Name:   name,
		Desc:   desc,
		Routes: make(map[string]Route),
		Quit:   make(chan bool),
	}
}

// A Connection connects two Routes
type Connection struct {
	from Route     // block Route to listen to
	to   Route     // block Route to send to
	quit chan bool // send a bool to this channel to remove the connection
}

// LastMessage returns the last message that flowed through connection
func (c *Connection) LastMessage() Msg {
	var out Msg
	return out
}

// Rate returns an estimate of the message flow rate in messages per second
func (c *Connection) Rate() float64 {
	return 1
}

// Serve handles passing messages between routes
func (c *Connection) Serve() {
	for {
		select {
		case m := <-c.from.C:
			c.to.C <- m
		case <-c.quit:
			return
		}
	}
}

// Stop closes the connection's outbound channel and causes the sends on the inbound Route to block
func (c *Connection) Stop() {
	close(c.to.C)
	c.quit <- true
}

// Connect builds a new connection and sets it running
func Connect(supervisor suture.Supervisor, from, to Route) {
	connection := &Connection{
		from: from,
		to:   to,
		quit: make(chan bool),
	}
	supervisor.Add(connection)
}

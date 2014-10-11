package main

import "github.com/thejerf/suture"

// Msg defines the individual messages that flow through streamtools
type Msg interface{}

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

type Connection struct {
	from Route     // block Route to listen to
	to   Route     // block Route to send to
	quit chan bool // send a bool to this channel to remove the connection
}

func (c *Connection) LastMessage() Msg {
	var out Msg
	return out
}

func (c *Connection) Rate() float64 {
	return 1
}

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

func (c *Connection) Stop() {
	close(c.to.C)
	c.quit <- true
}

func Connect(supervisor suture.Supervisor, from, to Route) {
	connection := &Connection{
		from: from,
		to:   to,
		quit: make(chan bool),
	}
	supervisor.Add(connection)
}

package main

import (
	"log"

	"github.com/thejerf/suture"
)

// Messages can be JSON, XML, GOJEE
type MsgType int

const (
	JSON MsgType = iota
	XML
	GOJEE
)

// Msg defines the individual messages that flow through streamtools
type Msg struct {
	Payload interface{}
	Type    MsgType
}

// Route controls inbound or outbound messages from a block
type Route struct {
	FromBlock      chan Msg
	ToBlock        chan Msg
	FromConnection chan Msg
	ToConnection   chan Msg
	Quit           chan bool // remove the Route from the block
}

func NewRoute() *Route {
	return &Route{
		FromBlock:      make(chan Msg),
		ToBlock:        make(chan Msg),
		FromConnection: make(chan Msg),
		ToConnection:   make(chan Msg),
		Quit:           make(chan bool),
	}
}

// starts accepting messages
func (r *Route) Start() {
	for {
		select {
		case m := <-r.FromBlock:
			r.ToConnection <- m
		case m := <-r.FromConnection:
			switch m.Type {
			case JSON:
				r.ToBlock <- m
			case GOJEE:
				log.Fatal("execute GOJEE before sending to block")
			case XML:
				log.Fatal("need parser")
			}
		case <-r.Quit:
			close(r.ToBlock)
			close(r.ToConnection)
		}
	}
}

// removes the route from its containg block
func (r *Route) Stop() {
	log.Println("removing Route")
	r.Quit <- true
}

// Block defines the basic properties of a streamtools block
type Block struct {
	Name   string           // friendly name used in the API and UI
	Desc   string           // short description of the block for the UI
	Routes map[string]Route // container for this block's routes
	Quit   chan bool        // channel to indicate this block should quit
}

func (b *Block) Serve() {

}

func (b *Block) Stop() {
	for _, route := range b.Routes {
		close(route.FromBlock)
		route.Quit <- true
	}
	b.Quit <- true
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
		case m := <-c.from.ToConnection:
			c.to.FromConnection <- m
		case <-c.quit:
			return
		}
	}
}

// Stop closes the connection's outbound channel and causes the sends on the inbound Route to block
func (c *Connection) Stop() {
	close(c.to.FromConnection)
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

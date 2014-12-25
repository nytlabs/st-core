package core

import (
	"sync"

	"github.com/nikhan/go-fetch"
)

const (
	NONE = iota
	KEY_VALUE
	PRIORITY
	ARRAY
	VALUE
)

type RouteID int
type SharedType int
type Message interface{}
type Connection chan Message
type MessageMap map[RouteID]Message

// Interrupt is a function that interrupts a running block in order to change its state.
// If the interrupt returns false, the block will quit.
type Interrupt func() bool

// Kernel is the core function that operates on an inbound message. It works by populating
// the outbound MessageMap, and can be interrupted on its Interrupt channel.
type Kernel func(MessageMap, MessageMap, Store, chan Interrupt) Interrupt

// A Pin is an inbound or outbound route to a block used in the Spec.
type Pin struct {
	Name string
}

// A Spec defines a block's input and output Pins, and the block's Kernel.
type Spec struct {
	Inputs  []Pin
	Outputs []Pin
	Shared  SharedType
	Kernel  Kernel
}

// A Route is an inbound route to a block. A Route holds the channel that allows Messages
// to be passed into the block. A Route's Path is applied to the inbound Message before populating the
// MessageMap and calling the Kernel. A Route can be set to a Value, instead of waiting for an inbound message.
type Route struct {
	Name  string
	Path  *fetch.Query
	Value *Message
	C     chan Message
}

// An Output holds a set of Connections. Each Connection refers to a Route.C. Every outbound
// mesage is sent on every Connection in the Connections set.
type Output struct {
	Name        string
	Connections map[Connection]struct{}
}

// A ManifestPair is a unique reference to an Output/Connection pair
type ManifestPair struct {
	string
	Connection
}

// A block's Manifest is the set of Connections
type Manifest map[ManifestPair]struct{}

// A block's BlockState is the pair of input/output MessageMaps, and the Manifest
type BlockState struct {
	inputValues  MessageMap
	outputValues MessageMap
	manifest     Manifest
	Processed    bool
}

type Store interface {
	Lock()
	Unlock()
}

type SharedStore struct {
	Type  SharedType
	Store Store
}

// A block's BlockRouting is the set of Input and Output routes, and the Interrupt channel
type BlockRouting struct {
	Inputs        []Route
	Outputs       []Output
	Shared        SharedStore
	InterruptChan chan Interrupt
	sync.RWMutex
}

// A Block is comprised of a state, a set of routes, and a kernel
type Block struct {
	state   BlockState
	routing BlockRouting
	kernel  Kernel
}

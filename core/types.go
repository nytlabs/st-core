package core

import "sync"

const (
	NONE = iota
	KEY_VALUE
	STREAM
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
type Kernel func(MessageMap, MessageMap, MessageMap, Store, chan Interrupt) Interrupt

// A Pin is an inbound or outbound route to a block used in the Spec.
type Pin struct {
	Name string
}

// A Spec defines a block's input and output Pins, and the block's Kernel.
type Spec struct {
	Name    string
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
	Value interface{}
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
	int
	Connection
}

// A block's Manifest is the set of Connections
type Manifest map[ManifestPair]struct{}

// A block's BlockState is the pair of input/output MessageMaps, and the Manifest
type BlockState struct {
	inputValues    MessageMap
	outputValues   MessageMap
	internalValues MessageMap
	manifest       Manifest
	Processed      bool
}

// a Store is esssentially a lockable piece of memory that can be accessed safely by mulitple blocks.
// The Lock and Unlock methods are usually implemented using a sync.Mutex
// TODO Store -> Source
type Store interface {
	Lock()
	Unlock()
}

// TODO collapse SharedStore into Store by blessing Store with a GetType() method
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

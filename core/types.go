package core

import (
	"errors"
	"sync"
	"time"
)

const (
	NONE = iota
	KEY_VALUE
	LIST
	VALUE_PRIMITIVE
	PRIORITY
	NSQCONSUMER
	WSCLIENT
	STDIN
)

// JSONType defines the possible types that variables in core can take
type JSONType uint8

const (
	NUMBER JSONType = iota
	STRING
	ARRAY
	OBJECT
	BOOLEAN
	NULL
	ANY
	WRITER
	ERROR
)

// BlockAlert defines the possible messages a block can emit about its runnig state
type BlockInfo uint8

const (
	BI_RUNNING BlockInfo = iota
	BI_ERROR
	BI_INPUT
	BI_OUTPUT
	BI_KERNEL
)

func (ba BlockInfo) MarshalJSON() ([]byte, error) {
	switch ba {
	case BI_RUNNING:
		return []byte(`"running"`), nil
	case BI_ERROR:
		return []byte(`"error"`), nil
	case BI_INPUT:
		return []byte(`"input"`), nil
	case BI_OUTPUT:
		return []byte(`"output"`), nil
	case BI_KERNEL:
		return []byte(`"kernel"`), nil
	}

	return nil, errors.New("could not marshal BlockAlert")
}

// MessageMap maps a block's inbound routes onto the Messages they contain
type MessageMap map[RouteIndex]Message

// Message is the container for data sent between blocks
type Message interface{}

// RouteIndex is the index into a MessageMap. The 0th index corresponds to that block's 0th Input or Output
type RouteIndex int

// SourceType is used to indicate what kind of source a block can connect to
type SourceType int

func (s *SourceType) UnmarshalJSON(data []byte) error {
	st := string(data)
	switch st {
	case `null`:
		*s = SourceType(NONE)
	case `"key_value"`:
		*s = SourceType(KEY_VALUE)
	case `"NSQ"`:
		*s = SourceType(NSQCONSUMER)
	case `"wsClient"`:
		*s = SourceType(WSCLIENT)
	case `"list"`:
		*s = SourceType(LIST)
	case `"value"`:
		*s = SourceType(VALUE_PRIMITIVE)
	case `"priority-queue"`:
		*s = SourceType(PRIORITY)
	case `"stdin"`:
		*s = SourceType(STDIN)
	default:
		return errors.New("Error unmarshalling source type")
	}
	return nil
}

func (s SourceType) MarshalJSON() ([]byte, error) {
	switch s {
	case NONE:
		return []byte(`null`), nil
	case KEY_VALUE:
		return []byte(`"key_value"`), nil
	case NSQCONSUMER:
		return []byte(`"NSQ"`), nil
	case WSCLIENT:
		return []byte(`"wsClient"`), nil
	case LIST:
		return []byte(`"list"`), nil
	case VALUE_PRIMITIVE:
		return []byte(`"value"`), nil
	case PRIORITY:
		return []byte(`"priority-queue"`), nil
	case STDIN:
		return []byte(`"stdin"`), nil
	}
	return nil, errors.New("Unknown source type")
}

// Connections are used to connect blocks together
type Connection chan Message

// Interrupt is a function that interrupts a running block in order to change its state.
// If the interrupt returns false, the block will quit.
type Interrupt func() bool

// Kernel is a block's core function that operates on an inbound message. It works by populating
// the outbound MessageMap, and can be interrupted on its Interrupt channel.
type Kernel func(MessageMap, MessageMap, MessageMap, Source, chan Interrupt) Interrupt

// A Pin contains information about a particular input or output
type Pin struct {
	Name string
	Type JSONType
}

// A Spec defines a block's input and output Pins, and the block's Kernel.
type Spec struct {
	Name     string
	Category []string
	Inputs   []Pin
	Outputs  []Pin
	Source   SourceType
	Kernel   Kernel
}

// Input is an inbound route to a block. A Input holds the channel that allows Messages
// to be passed into the block. A Input's Path is applied to the inbound Message before populating the
// MessageMap and calling the Kernel. A Input can be set to a Value, instead of waiting for an inbound message.
type Input struct {
	Name  string       `json:"name"`
	Value *InputValue  `json:"value"`
	Type  JSONType     `json:"type"`
	C     chan Message `json:"-"`
}

type InputValue struct {
	Data interface{} `json:"data"`
}

func (i *InputValue) Exists() bool {
	return i != nil
}

// An Output holds a set of Connections. Each Connection refers to a Input.C. Every outbound
// mesage is sent on every Connection in the Connections set.
type Output struct {
	Name        string                  `json:"name"`
	Type        JSONType                `json:"type"`
	Connections map[Connection]struct{} `json:"-"`
}

// A SourceSpec defines a source's name and type
type SourceSpec struct {
	Name     string
	Type     SourceType
	New      SourceFunc
	Category []string
}

// A function that creates a source
type SourceFunc func() Source

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

type Source interface {
	GetType() SourceType
}

type Interface interface {
	Source
	Serve()
	Stop()
}

type Store interface {
	Source
	Get() interface{}
	Set(interface{}) error
	Lock()
	Unlock()
}

// A block's BlockRouting is the set of Input and Output routes, and the Interrupt channel
type BlockRouting struct {
	Inputs        []Input
	Outputs       []Output
	Source        Source
	InterruptChan chan Interrupt
	sync.RWMutex
}

// A Block describes the block's components
type Block struct {
	state      BlockState
	routing    BlockRouting
	kernel     Kernel
	sourceType SourceType
	Monitor    chan MonitorMessage
	lastCrank  time.Time
	done       chan struct{}
	//blockageTimer *time.Timer
}

type MonitorMessage struct {
	Type BlockInfo   `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

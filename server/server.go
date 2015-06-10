package server

import (
	"log"
	"sync"

	"github.com/nytlabs/st-core/core"
)

const (
	// actions
	DELETE = "delete"
	RESET  = "reset"
	UPDATE = "update"
	CREATE = "create"
	ALERT  = "alert"
	// nodes
	BLOCK  = "block"
	GROUP  = "group"
	SOURCE = "source"
	// edges
	LINK       = "link"
	CONNECTION = "connection"
	// attributes
	CHILD = "child"
	ROUTE = "route"
	PARAM = "param"
)

// The Server maintains a set of handlers that coordinate the creation of Nodes
type Server struct {
	groups        map[int]*Group // TODO these maps aren't strictly necessary, but save constantly performing depth first searches
	parents       map[int]int
	blocks        map[int]*BlockLedger
	connections   map[int]*ConnectionLedger
	sources       map[int]*SourceLedger
	links         map[int]*LinkLedger
	library       map[string]core.Spec
	sourceLibrary map[string]core.SourceSpec
	lastID        int
	addSocket     chan *socket
	delSocket     chan *socket
	broadcast     chan []byte
	emitChan      chan []byte
	sync.Mutex
}

// NewServer starts a new Server. This object is immediately up and running.
func NewServer() *Server {
	groups := make(map[int]*Group)
	groups[0] = &Group{
		Label:    "root",
		Id:       0,
		Children: []int{},
		Parent:   nil,
	}

	blocks := make(map[int]*BlockLedger)
	connections := make(map[int]*ConnectionLedger)
	sources := make(map[int]*SourceLedger)
	links := make(map[int]*LinkLedger)
	library := core.GetLibrary()
	sourceLibrary := core.GetSources()
	parents := make(map[int]int)
	s := &Server{
		lastID:        0,
		parents:       parents,
		groups:        groups,
		blocks:        blocks,
		sourceLibrary: sourceLibrary,
		connections:   connections,
		library:       library,
		links:         links,
		sources:       sources,
		addSocket:     make(chan *socket),
		delSocket:     make(chan *socket),
		broadcast:     make(chan []byte),
		emitChan:      make(chan []byte),
	}
	// ws stuff
	log.Println("starting websocker handler")
	go s.websocketRouter()
	return s
}

// GetNextID eturns the next ID to be used for a new group or a new block
func (s *Server) GetNextID() int {
	s.lastID++
	return s.lastID
}

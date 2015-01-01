package server

import (
	"log"
	"sync"

	"github.com/nytlabs/st-core/core"
	"github.com/thejerf/suture"
)

const (
	DELETE     = "DELETE"
	MODIFY     = "MODIFY"
	CREATE     = "CREATE"
	CONNECTION = "CONNECTION"
	BLOCK      = "BLOCK"
	GROUP      = "GROUP"
)

// The Server maintains a set of handlers that coordinate the creation of Nodes
type Server struct {
	groups      map[int]*Group // TODO these maps aren't strictly necessary, but save constantly performing depth first searches
	blocks      map[int]*BlockLedger
	connections map[int]*ConnectionLedger
	stores      map[int]*core.Store
	library     map[string]core.Spec
	supervisor  *suture.Supervisor
	lastID      int
	addSocket   chan *socket
	delSocket   chan *socket
	broadcast   chan []byte
	emitChan    chan []byte
	sync.Mutex
}

// NewServer starts a new Server. This object is immediately up and running.
func NewServer() *Server {
	supervisor := suture.NewSimple("st-core")
	supervisor.ServeBackground()
	groups := make(map[int]*Group)
	groups[0] = NewGroup(0, "root") // this is the top level group
	blocks := make(map[int]*BlockLedger)
	connections := make(map[int]*ConnectionLedger)
	stores := make(map[int]*core.Store)
	library := core.GetLibrary()
	s := &Server{
		supervisor:  supervisor,
		lastID:      0,
		groups:      groups,
		blocks:      blocks,
		connections: connections,
		library:     library,
		stores:      stores,
		addSocket:   make(chan *socket),
		delSocket:   make(chan *socket),
		broadcast:   make(chan []byte),
		emitChan:    make(chan []byte),
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

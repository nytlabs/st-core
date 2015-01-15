package server

import (
	"log"
	"sync"

	"github.com/nytlabs/st-core/core"
	"github.com/thejerf/suture"
)

const (
	DELETE      = "delete"
	UPDATE      = "update"
	CREATE      = "create"
	CONNECTION  = "connection"
	BLOCK       = "block"
	GROUP       = "group"
	GROUP_CHILD = "group_child"
)

// The Server maintains a set of handlers that coordinate the creation of Nodes
type Server struct {
	groups      map[int]*Group // TODO these maps aren't strictly necessary, but save constantly performing depth first searches
	parents     map[int]int
	blocks      map[int]*BlockLedger
	connections map[int]*ConnectionLedger
	//routes      map[RoutePair]RoutePair
	stores     map[int]*core.Store
	library    map[string]core.Spec
	supervisor *suture.Supervisor
	lastID     int
	addSocket  chan *socket
	delSocket  chan *socket
	broadcast  chan []byte
	emitChan   chan []byte
	sync.Mutex
}

// NewServer starts a new Server. This object is immediately up and running.
func NewServer() *Server {
	supervisor := suture.NewSimple("st-core")
	supervisor.ServeBackground()
	groups := make(map[int]*Group)
	//groups[0] = NewGroup(0, "root") // this is the top level group
	groups[0] = &Group{
		Label:    "root",
		Id:       0,
		Children: []int{},
		Parent:   nil,
	}

	blocks := make(map[int]*BlockLedger)
	connections := make(map[int]*ConnectionLedger)
	stores := make(map[int]*core.Store)
	library := core.GetLibrary()
	parents := make(map[int]int)
	//routes := make(map[RoutePair]RoutePair)
	s := &Server{
		supervisor:  supervisor,
		lastID:      0,
		parents:     parents,
		groups:      groups,
		blocks:      blocks,
		connections: connections,
		library:     library,
		stores:      stores,
		//routes:      routes,
		addSocket: make(chan *socket),
		delSocket: make(chan *socket),
		broadcast: make(chan []byte),
		emitChan:  make(chan []byte),
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

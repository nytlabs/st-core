package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/thejerf/suture"
)

// MOCK Block

type Block struct {
	q     chan bool
	token suture.ServiceToken
	id    int
	name  string
}

type BlockJSON struct {
	Name string
}

func (b *Block) Serve() {
	<-b.q
}

func (b *Block) Stop() {
	b.q <- true
}

func (b *Block) UnmarshalJSON(msg []byte) error {
	var m BlockJSON
	err := json.Unmarshal(msg, &m)
	if err != nil {
		log.Panic(err)
	}
	b.name = m.Name
	return nil
}

func (b *Block) GetID() int {
	return b.id
}

func NewBlock(id int) *Block {
	q := make(chan bool)
	b := &Block{
		id: id,
		q:  q,
	}
	return b
}

// A Node is simply something that can report its ID
type Node interface {
	GetID() int
}

// A Group is contained by a node and contains nodes
type Group struct {
	children map[Node]bool
	parent   Node
	id       int
	sync.Mutex
}

func (g *Group) GetID() int {
	return g.id
}

// NewGroup returns a group with no children
func NewGroup(id int, parent Node) *Group {
	return &Group{
		children: make(map[Node]bool),
		id:       id,
		parent:   parent,
	}
}

// NewGroupFromNodes returns a group containing existing nodes
func NewGroupFromNodes(id int, parent Node, nodes []Node) *Group {
	g := NewGroup(id, parent)
	for _, n := range nodes {
		g.AddNode(n)
	}
	return g
}

// AddNode adds a node to a group
func (g *Group) AddNode(n Node) {
	g.Lock()
	defer g.Unlock()
	g.children[n] = true
}

// RemoveNode removes a node from a group
func (g *Group) RemoveNode(n Node) {
	g.Lock()
	defer g.Unlock()
	delete(g.children, n)
}

// The Server maintains a set of handlers that coordinate the creation of Nodes
type Server struct {
	groups     map[int]*Group
	supervisor *suture.Supervisor
	lastID     int
	sync.Mutex
}

func NewServer() *Server {
	supervisor := suture.NewSimple("st-core")
	supervisor.ServeBackground()
	groups := make(map[int]*Group)
	groups[0] = NewGroup(0, nil) // this is the top level group
	return &Server{
		supervisor: supervisor,
		lastID:     0,
		groups:     groups,
	}
}

func (s *Server) GetNextID() int {
	s.lastID += 1
	return s.lastID
}

func (s *Server) createBlockHandler(w http.ResponseWriter, r *http.Request) {
	b := NewBlock(s.GetNextID())
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(body)
	}
	log.Println("adding", b.name, "to group 0 with id", b.id)
	// we need to introduce the block to our running supervisor
	b.token = s.supervisor.Add(b)
	// and we need to assign it to a group
	s.groups[0].AddNode(b)
}

func main() {
	s := NewServer()
	r := mux.NewRouter()
	//r.HandleFunc("/", s.rootHandler).Methods("GET")
	//r.HandleFunc("/group", s.createGroupHandler).Methods("POST")
	r.HandleFunc("/block", s.createBlockHandler).Methods("POST")
	http.Handle("/", r)

	log.Println("serving on 7071")
	err := http.ListenAndServe(":7071", nil)
	if err != nil {
		log.Panicf(err.Error())
	}
}

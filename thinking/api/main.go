package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	children map[int]Node
	id       int
	sync.Mutex
}

// GetID returns the id of the Group
func (g *Group) GetID() int {
	return g.id
}

// NewGroup returns a group with no children
func NewGroup(id int) *Group {
	return &Group{
		children: make(map[int]Node),
		id:       id,
	}
}

// NewGroupFromNodes returns a group containing existing nodes
func NewGroupFromNodes(id int, nodes []Node) *Group {
	g := NewGroup(id)
	for _, n := range nodes {
		g.AddNode(n)
	}
	return g
}

// GetNode returns the specified Node, erroring if the node does not exist
func (g *Group) GetNode(id int) (Node, error) {
	out, ok := g.children[id]
	if !ok {
		return nil, errors.New("could not find node with id " + strconv.Itoa(id))
	}
	return out, nil
}

// AddNode adds a node to a group
func (g *Group) AddNode(n Node) {
	g.Lock()
	defer g.Unlock()
	g.children[n.GetID()] = n
}

// RemoveNode removes a node from a group
func (g *Group) RemoveNode(n Node) {
	g.Lock()
	defer g.Unlock()
	delete(g.children, n.GetID())
}

// The Server maintains a set of handlers that coordinate the creation of Nodes
type Server struct {
	groups     map[int]*Group
	supervisor *suture.Supervisor
	lastID     int
	sync.Mutex
}

// Starts a new Server. This object is immediately up and running.
func NewServer() *Server {
	supervisor := suture.NewSimple("st-core")
	supervisor.ServeBackground()
	groups := make(map[int]*Group)
	groups[0] = NewGroup(0) // this is the top level group
	return &Server{
		supervisor: supervisor,
		lastID:     0,
		groups:     groups,
	}
}

// Returns the next ID to be used for a new group or a new block
func (s *Server) GetNextID() int {
	s.lastID += 1
	return s.lastID
}

// Creates a new block and adds it to the Server.
// TODO currently all blocks start off life in the root group, which may be a bit limiting.
func (s *Server) createBlockHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(body)
	}
	b := NewBlock(s.GetNextID())
	err = json.Unmarshal(body, &b)
	if err != nil {
		log.Panic(err)
	}
	log.Println("adding", b.name, "to group 0 with id", b.id)
	// we need to introduce the block to our running supervisor
	b.token = s.supervisor.Add(b)
	// and we need to assign it to a group
	s.groups[0].AddNode(b)
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

// when creating a new group you must specify the parent of the new group and the children (blocks and other groups) of the new group.
type createGroupRequest struct {
	ParentID int
	ChildIDs []int
}

// creates a new group and adds it to the Server. Moves all of the specified children out of the parent's group and into the new group.
func (s *Server) createGroupHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(body)
	}
	var groupReq createGroupRequest
	json.Unmarshal(body, &groupReq)

	g := NewGroup(s.GetNextID())
	if err != nil {
		log.Panic(err)
	}

	// move over the children into the new group
	parent := s.groups[groupReq.ParentID] // TODO get this by recursing through the tree
	for _, childid := range groupReq.ChildIDs {
		child, err := parent.GetNode(childid)
		if err != nil {
			log.Panic(err)
		}
		log.Println("removing", child.GetID(), "from group", parent.id)
		parent.RemoveNode(child)
		g.AddNode(child)
	}

	log.Println("adding new group as child of", groupReq.ParentID, "with ID", g.id)
	parent.AddNode(g)
	s.groups[g.id] = g

	// reassurance
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func printGroups(g *Group, out string, tabs string) (string, string) {
	tabs += "  "
	for _, child := range g.children {
		out += tabs + strconv.Itoa(child.GetID())
		switch child := child.(type) {
		case *Group:
			out += "-\n"
			out, tabs = printGroups(child, out, tabs)
		case *Block:
			out += "\n"
		}
	}
	tabs = tabs[len(tabs)-2 : len(tabs)]
	return out, tabs
}

func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	out, _ := printGroups(s.groups[0], "0 -\n", "")
	w.WriteHeader(200)
	w.Write([]byte(out))
}

func main() {
	s := NewServer()
	r := mux.NewRouter()
	r.HandleFunc("/", s.rootHandler).Methods("GET")
	r.HandleFunc("/group", s.createGroupHandler).Methods("POST")
	r.HandleFunc("/block", s.createBlockHandler).Methods("POST")
	http.Handle("/", r)

	log.Println("serving on 7071")
	err := http.ListenAndServe(":7071", nil)
	if err != nil {
		log.Panicf(err.Error())
	}
}

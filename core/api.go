package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/thejerf/suture"
)

// BlockLedger is the information the API keeps about the block
// TODO Token and ID should probably be the same thing
type BlockLedger struct {
	Name      string
	BlockType string
	Token     suture.ServiceToken
	ID        int
	Block     *Block
}

// GetID returns the block's ID
func (b BlockLedger) GetID() int {
	return b.ID
}

// GetName returns the block's Name
func (b BlockLedger) GetName() string {
	return b.Name
}

// A Node has an ID and a Name
type Node interface {
	GetID() int
	GetName() string
}

// A Group is contained by a node and contains nodes
type Group struct {
	children map[int]Node
	id       int
	name     string
	sync.Mutex
}

// GetID returns the id of the Group
func (g *Group) GetID() int {
	return g.id
}

// GetName returns the group's Name
func (g *Group) GetName() string {
	return g.name
}

// NewGroup returns a group with no children
func NewGroup(id int, name string) *Group {
	return &Group{
		children: make(map[int]Node),
		id:       id,
		name:     name,
	}
}

// NewGroupFromNodes returns a group containing existing nodes
func NewGroupFromNodes(id int, name string, nodes []Node) *Group {
	g := NewGroup(id, name)
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
	groups     map[int]*Group // TODO these maps aren't strictly necessary, but save constantly performing depth first searches
	blocks     map[int]*BlockLedger
	supervisor *suture.Supervisor
	lastID     int
	sync.Mutex
}

// NewServer starts a new Server. This object is immediately up and running.
func NewServer() *Server {
	supervisor := suture.NewSimple("st-core")
	supervisor.ServeBackground()
	groups := make(map[int]*Group)
	groups[0] = NewGroup(0, "root") // this is the top level group
	blocks := make(map[int]*BlockLedger)
	return &Server{
		supervisor: supervisor,
		lastID:     0,
		groups:     groups,
		blocks:     blocks,
	}
}

// GetNextID eturns the next ID to be used for a new group or a new block
func (s *Server) GetNextID() int {
	s.lastID++
	return s.lastID
}

// CreateBlockHandler responds to a POST request to instantiate a new block and add it to the Server.
// TODO currently all blocks start off life in the root group, which may be a bit limiting.
func (s *Server) CreateBlockHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(body)
	}
	var m BlockLedger
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Panic(err)
	}
	blockSpec, ok := Library[m.BlockType]
	if !ok {
		log.Fatal("unknown block type")
	}
	b := NewBlock(blockSpec)
	m.ID = s.GetNextID()
	m.Block = b
	log.Println("adding", m.Name, "to group 0 with id", m.ID)
	// we need to introduce the block to our running supervisor
	m.Token = s.supervisor.Add(b)
	// and we need to assign it to a group
	s.groups[0].AddNode(&m)
	s.blocks[m.ID] = &m
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

// when creating a new group you must specify the parent of the new group and the children (blocks and other groups) of the new group.
type createGroupRequest struct {
	ParentID int
	Name     string
	ChildIDs []int
}

// CreateGroupHandler responds to a POST request to instantiate a new group and add it to the Server.
// Moves all of the specified children out of the parent's group and into the new group.
func (s *Server) CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(body)
	}
	var groupReq createGroupRequest
	json.Unmarshal(body, &groupReq)

	g := NewGroup(s.GetNextID(), groupReq.Name)
	if err != nil {
		log.Panic(err)
	}

	// move over the children into the new group
	parent := s.groups[groupReq.ParentID]
	for _, childid := range groupReq.ChildIDs {
		child, err := parent.GetNode(childid)
		if err != nil {
			log.Panic(err)
		}
		log.Println("removing node", child.GetID(), "from group", parent.GetID())
		parent.RemoveNode(child)
		log.Println("adding node", child.GetID(), "to group", g.GetID())
		g.AddNode(child)
	}

	log.Println("adding new group", g.GetName(), "as child of", groupReq.ParentID, "with ID", g.GetID())
	parent.AddNode(g)
	s.groups[g.id] = g

	// reassurance
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

type connectRequest struct {
	FromBlockID    int
	FromRouteIndex int
	ToBlockID      int
	ToRouteIndex   int
}

// CreateConnectionHandler responds to a POST request to instantiate a new connection
func (s *Server) CreateConnectionHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(body)
	}
	var connectReq connectRequest
	json.Unmarshal(body, &connectReq)

	from := s.blocks[connectReq.FromBlockID]
	to := s.blocks[connectReq.ToBlockID]
	fromRoute := RouteID(connectReq.FromRouteIndex)
	toRoute := to.Block.Input(RouteID(connectReq.ToRouteIndex))

	log.Printf("Connecting %v:%v -> %v:%v\n", from.Name, connectReq.FromRouteIndex, connectReq.ToRouteIndex, to.Name)
	from.Block.Connect(fromRoute, toRoute.C)

	// reassurance
	w.WriteHeader(200)
	w.Write([]byte("OK"))

}

func printGroups(g *Group, out string, tabs string) (string, string) {
	tabs += "  "
	for _, child := range g.children {
		switch child := child.(type) {
		case *Group:
			out += tabs + "(" + strconv.Itoa(child.GetID()) + ") " + child.GetName() + ":\n"
			out, tabs = printGroups(child, out, tabs)
		case *BlockLedger:
			out += tabs + "(" + strconv.Itoa(child.GetID()) + ") " + child.GetName() + "\n"
		default:
			log.Printf("%T", child)
			log.Fatal("trying to print some shit")
		}
	}
	tabs = tabs[len(tabs)-2 : len(tabs)]
	return out, tabs
}

func (s *Server) String() string {
	out, _ := printGroups(s.groups[0], "(root):\n", "")
	return out
}

// RootHandler returns a string representation of the groups in the Server. This won't last!
func (s *Server) RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprint(s)))
}

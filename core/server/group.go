package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// A Node has an ID and a Name
type Node interface {
	GetID() int
	GetName() string
	json.Marshaler
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

// MarshalJSON returns the JSON representation of the group
func (g *Group) MarshalJSON() ([]byte, error) {
	children := make([]Node, len(g.children))
	i := 0
	for _, v := range g.children {
		children[i] = v
		i++
	}
	out := map[string]interface{}{
		"children": children,
		"id":       strconv.Itoa(g.id),
		"name":     g.name,
	}
	return json.Marshal(out)
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

// when creating a new group you must specify the parent of the new group and the children (blocks and other groups) of the new group.
type createGroupRequest struct {
	ParentID int
	Name     string
	ChildIDs []int
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

// GetGroupHandler returns a string representation of the groups in the Server. This won't last!
func (s *Server) GetGroupHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprint(s)))
}

func (s *Server) GroupIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(struct{}{}); err != nil {
		panic(err)
	}
}

// CreateGroupHandler responds to a POST request to instantiate a new group and add it to the Server.
// Moves all of the specified children out of the parent's group and into the new group.
func (s *Server) GroupCreate(w http.ResponseWriter, r *http.Request) {
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

	// broadcast state update
	log.Println(g)
	stateUpdate, err := json.Marshal(g)
	if err != nil {
		log.Fatal(err)
	}
	s.broadcast <- stateUpdate

	// reassurance
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
func (s *Server) GroupDelete(w http.ResponseWriter, r *http.Request) {
}

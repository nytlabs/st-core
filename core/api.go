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
	"time"

	"github.com/gorilla/websocket"
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

// MarshalJSON returns the JSON representation of the block ledger
func (b BlockLedger) MarshalJSON() ([]byte, error) {
	out := map[string]interface{}{
		"id":        strconv.Itoa(b.ID),
		"name":      b.Name,
		"blockType": b.BlockType,
	}
	return json.Marshal(out)
}

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

// The Server maintains a set of handlers that coordinate the creation of Nodes
type Server struct {
	groups     map[int]*Group // TODO these maps aren't strictly necessary, but save constantly performing depth first searches
	blocks     map[int]*BlockLedger
	stores     map[int]*StateLocker
	library    map[string]Spec
	supervisor *suture.Supervisor
	lastID     int
	addConn    chan *connection
	delConn    chan *connection
	broadcast  chan []byte
	emitChan   chan []byte
	sync.Mutex
}

// NewServer starts a new Server. This object is immediately up and running.
func NewServer() *Server {
	supervisor := suture.NewSimple("st-core")
	supervisor.ServeBackground()
	groups := make(map[int]*Group)
	groups[0] = NewGroup(0, "root") // this is the top level group
	blocks := make(map[int]*BlockLedger)
	stores := make(map[int]*StateLocker)
	library := GetLibrary()
	s := &Server{
		supervisor: supervisor,
		lastID:     0,
		groups:     groups,
		blocks:     blocks,
		library:    library,
		stores:     stores,
		addConn:    make(chan *connection),
		delConn:    make(chan *connection),
		broadcast:  make(chan []byte),
		emitChan:   make(chan []byte),
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

// WEBSOCKET STUFF

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type connection struct {
	ws   *websocket.Conn
	send chan []byte
}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (s *Server) websocketRouter() {
	hub := make(map[*connection]bool)
	for {
		select {
		case c := <-s.addConn:
			hub[c] = true
		case c := <-s.delConn:
			delete(hub, c)
		case m := <-s.broadcast:
			for c := range hub {
				c.send <- m
			}
		}
	}
}

func (s *Server) websocketReadPump(c *connection) {
	defer func() {
		s.delConn <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		s.emitChan <- message
	}
}

func (s *Server) websocketWritePump(c *connection) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// WebsocketHandler upgrades the inbound request to a websocket, establashing the necessary machinery to allow
// two way communication between st-core and the client
func (s *Server) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	s.addConn <- c
	go s.websocketWritePump(c)
	go s.websocketReadPump(c)
}

func (s *Server) addHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"server":"ok"}`))
}

// HTTP HANDLERS

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
	blockSpec, ok := s.library[m.BlockType]
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
	//broadcast for state update
	log.Println("broadcasting")
	s.broadcast <- body
	log.Println("done")
	// write HTTP response
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

	// broadcast state update
	s.broadcast <- body

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

// GetGroupHandler returns a string representation of the groups in the Server. This won't last!
func (s *Server) GetGroupHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprint(s)))
}

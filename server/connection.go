package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nytlabs/st-core/core"
)

type ConnectionNode struct {
	Id    int `json:"id"`
	Route int `json:"route"`
}

type ConnectionLedger struct {
	Source ConnectionNode `json:"from"`
	Target ConnectionNode `json:"to"`
	Id     int            `json:"id"`
}

type ProtoConnection struct {
	Source ConnectionNode `json:"from"`
	Target ConnectionNode `json:"to"`
}

func (s *Server) ListConnections() []ConnectionLedger {
	connections := []ConnectionLedger{}
	for _, c := range s.connections {
		connections = append(connections, *c)
	}
	return connections
}

func (s *Server) ConnectionIndexHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()

	c := s.ListConnections()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(c); err != nil {
		panic(err)
	}
}

// CreateConnectionHandler responds to a POST request to instantiate a new connection
func (s *Server) ConnectionCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	var newConn ProtoConnection
	json.Unmarshal(body, &newConn)

	s.Lock()
	defer s.Unlock()

	nc, err := s.CreateConnection(newConn)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	writeJSON(w, nc)
}

func (s *Server) CreateConnection(newConn ProtoConnection) (*ConnectionLedger, error) {
	source, ok := s.blocks[newConn.Source.Id]
	if !ok {
		return nil, errors.New("source block does not exist")
	}

	target, ok := s.blocks[newConn.Target.Id]
	if !ok {
		return nil, errors.New("target block does not exist")
	}

	sourceRoute := core.RouteIndex(newConn.Source.Route)
	targetRoute, err := target.Block.GetInput(core.RouteIndex(newConn.Target.Route))
	if err != nil {
		return nil, err
	}

	err = source.Block.Connect(sourceRoute, targetRoute.C)
	if err != nil {
		return nil, err
	}

	conn := &ConnectionLedger{
		Source: newConn.Source,
		Target: newConn.Target,
		Id:     s.GetNextID(),
	}

	s.ResetGraph(conn)

	s.connections[conn.Id] = conn

	s.websocketBroadcast(Update{Action: CREATE, Type: CONNECTION, Data: wsConnection{*conn}})
	return conn, nil
}

/*
ResetGraph stops/resets/starts the entire connected subgraph related to
connection Conn. It is a general approach that probably touches a lot more
blocks than it needs to.

TODO: This function should be replaced with language-level
features that allow insight into stuck messages OR run in a more precise
fashion (detecting branching/merging of streams).
*/
func (s *Server) ResetGraph(conn *ConnectionLedger) {
	found := make(map[int]struct{})
	cacheST := make(map[int]map[int]struct{})
	cacheTS := make(map[int]map[int]struct{})

	for _, c := range s.connections {
		if _, ok := cacheST[c.Source.Id]; !ok {
			cacheST[c.Source.Id] = make(map[int]struct{})
		}
		if _, ok := cacheTS[c.Target.Id]; !ok {
			cacheTS[c.Target.Id] = make(map[int]struct{})
		}
		cacheST[c.Source.Id][c.Target.Id] = struct{}{}
		cacheTS[c.Target.Id][c.Source.Id] = struct{}{}
	}

	if _, ok := cacheST[conn.Source.Id]; !ok {
		cacheST[conn.Source.Id] = make(map[int]struct{})
	}

	if _, ok := cacheTS[conn.Target.Id]; !ok {
		cacheTS[conn.Target.Id] = make(map[int]struct{})
	}

	cacheST[conn.Source.Id][conn.Target.Id] = struct{}{}
	cacheTS[conn.Target.Id][conn.Source.Id] = struct{}{}

	var traverse func(int)

	// make a set of all nodes connecting to this connection
	traverse = func(id int) {
		found[id] = struct{}{}
		if _, ok := cacheST[id]; ok {
			for k, _ := range cacheST[id] {
				if _, ok = found[k]; !ok {
					traverse(k)
				}
			}
		}
		if _, ok := cacheTS[id]; ok {
			for k, _ := range cacheTS[id] {
				if _, ok = found[k]; !ok {
					traverse(k)
				}
			}
		}
	}

	traverse(conn.Source.Id)

	for k, _ := range found {
		log.Println("tidy: stopping id", k, s.blocks[k].Type)
		s.blocks[k].Block.Stop()
	}

	for k, _ := range found {
		log.Println("tidy: resetting id", k)
		s.blocks[k].Block.Reset()
	}

	for k, _ := range found {
		log.Println("tidy: starting id", k)
		go s.blocks[k].Block.Serve()
	}
}

// returns a description of the connection
func (s *Server) ConnectionHandler(w http.ResponseWriter, r *http.Request) {

	id, err := getIDFromMux(mux.Vars(r))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, err)
		return
	}

	conn, ok := s.connections[id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"could not find connection" + string(id)})
		return
	}

	w.WriteHeader(http.StatusOK)
	writeJSON(w, conn)
}

func (s *Server) ConnectionModifyCoordinates(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) DeleteConnection(id int) error {
	c, ok := s.connections[id]
	if !ok {
		return errors.New("could not find connection")
	}

	source, ok := s.blocks[c.Source.Id]
	if !ok {
		return errors.New("could not find source block")
	}

	target, ok := s.blocks[c.Target.Id]
	if !ok {
		return errors.New("could not find target block")
	}

	route, err := target.Block.GetInput(core.RouteIndex(c.Target.Route))
	if err != nil {
		return err
	}

	err = source.Block.Disconnect(core.RouteIndex(c.Source.Route), route.C)
	if err != nil {
		return err
	}

	delete(s.connections, id)

	s.ResetGraph(c)

	s.websocketBroadcast(Update{Action: DELETE, Type: CONNECTION, Data: wsConnection{wsId{id}}})
	return nil
}

func (s *Server) ConnectionDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromMux(mux.Vars(r))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, err)
		return
	}

	s.Lock()
	defer s.Unlock()

	err = s.DeleteConnection(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

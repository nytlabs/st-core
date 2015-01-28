package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nytlabs/st-core/core"
)

type ConnectionNode struct {
	Id    int `json:"id"`
	Route int `json:"route"`
}

type ConnectionLedger struct {
	Source ConnectionNode `json:"source"`
	Target ConnectionNode `json:"target"`
	Id     int            `json:"id"`
}

type ProtoConnection struct {
	Source ConnectionNode `json:"source"`
	Target ConnectionNode `json:"target"`
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

	s.connections[conn.Id] = conn

	s.websocketBroadcast(Update{Action: CREATE, Type: CONNECTION, Data: conn})
	return conn, nil
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

	s.websocketBroadcast(Update{Action: DELETE, Type: CONNECTION, Data: c})
	return nil
}

func (s *Server) ConnectionDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ids, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"no ID supplied"})
		return
	}

	id, err := strconv.Atoi(ids)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
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

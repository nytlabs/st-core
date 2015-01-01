package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nytlabs/st-core/core"
)

type ConnectionLedger struct {
	FromBlockID    int
	FromRouteIndex int
	ToBlockID      int
	ToRouteIndex   int
	ID             int
}

func (s *Server) ListConnections() []ConnectionLedger {
	var connections []ConnectionLedger
	s.Lock()
	for _, c := range s.connections {
		connections = append(connections, *c)
	}
	s.Unlock()
	return connections
}

func (s *Server) ConnectionIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(s.ListConnections()); err != nil {
		panic(err)
	}
}

// CreateConnectionHandler responds to a POST request to instantiate a new connection
func (s *Server) ConnectionCreate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(body)
	}
	var connectReq ConnectionLedger
	json.Unmarshal(body, &connectReq)

	from := s.blocks[connectReq.FromBlockID]
	to := s.blocks[connectReq.ToBlockID]
	fromRoute := core.RouteID(connectReq.FromRouteIndex)
	toRoute, err := to.Block.GetRoute(core.RouteID(connectReq.ToRouteIndex))
	if err != nil {
		log.Println("error:", err)
	}

	from.Block.Connect(fromRoute, toRoute.C)

	connectReq.ID = s.GetNextID()
	s.connections[connectReq.ID] = &connectReq

	out, _ := json.Marshal(connectReq)

	// broadcast state update
	s.broadcast <- out

	// reassurance
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
func (s *Server) ConnectionDelete(w http.ResponseWriter, r *http.Request) {
}

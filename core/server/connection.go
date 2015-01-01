package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nytlabs/st-core/core"
)

type connectRequest struct {
	FromBlockID    int
	FromRouteIndex int
	ToBlockID      int
	ToRouteIndex   int
}

func (s *Server) ConnectionIndex(w http.ResponseWriter, r *http.Request) {
}

// CreateConnectionHandler responds to a POST request to instantiate a new connection
func (s *Server) ConnectionCreate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(body)
	}
	var connectReq connectRequest
	json.Unmarshal(body, &connectReq)

	from := s.blocks[connectReq.FromBlockID]
	to := s.blocks[connectReq.ToBlockID]
	fromRoute := core.RouteID(connectReq.FromRouteIndex)
	toRoute, err := to.Block.GetRoute(core.RouteID(connectReq.ToRouteIndex))
	if err != nil {
		log.Println("error:", err)
	}

	log.Printf("Connecting %v:%v -> %v:%v\n", from.Name, connectReq.FromRouteIndex, connectReq.ToRouteIndex, to.Name)
	from.Block.Connect(fromRoute, toRoute.C)

	// broadcast state update
	s.broadcast <- body

	// reassurance
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
func (s *Server) ConnectionDelete(w http.ResponseWriter, r *http.Request) {
}

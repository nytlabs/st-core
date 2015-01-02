package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nytlabs/st-core/core"
	"github.com/thejerf/suture"
)

type BlockLedger struct {
	Name        string              `json:"name"`
	Id          int                 `json:"id"`
	Block       *core.Block         `json:"-"`
	Token       suture.ServiceToken `json:"-"`
	Parent      int                 `json:"parent"`
	Composition int                 `json:"composition,omitempty"`
	Spec        core.Spec           `json:"spec"`
}

type BlockValue struct {
	Id    int                    `json:"id"`
	Value map[string]interface{} `json:"value"`
}

func (s *Server) ListBlocks() []BlockLedger {
	var blocks []BlockLedger
	s.Lock()
	for _, b := range s.blocks {
		blocks = append(blocks, *b)
	}
	s.Unlock()
	return blocks
}

func (s *Server) BlockIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(s.ListBlocks()); err != nil {
		panic(err)
	}
}

// CreateBlockHandler responds to a POST request to instantiate a new block and add it to the Server.
func (s *Server) BlockCreate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"could not read request body"})
		return
	}

	var m BlockLedger
	err = json.Unmarshal(body, &m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"no ID supplied"})
		return
	}

	blockSpec, ok := s.library[m.Spec.Name]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"block not found"})
		return
	}

	m.Spec = blockSpec
	m.Id = s.GetNextID()
	m.Block = core.NewBlock(blockSpec)
	m.Token = s.supervisor.Add(m.Block)
	s.blocks[m.Id] = &m

	s.websocketBroadcast(Update{Action: CREATE, Type: BLOCK, Data: m})
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) BlockDelete(w http.ResponseWriter, r *http.Request) {
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

	b, ok := s.blocks[id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"connection not found"})
		return
	}

	deleteSet := make(map[int]struct{})

	// build a set of connections that we may need to delete
	// we need to panic here because if any error is thrown we are in huge trouble
	// any panic indicates that our server connection ledger is no longer true
	for _, c := range s.connections {
		if c.Target.Id == id {
			route, err := b.Block.GetRoute(core.RouteID(c.Target.Route))
			if err != nil {
				panic(err)
			}
			err = s.blocks[c.Source.Id].Block.Disconnect(core.RouteID(c.Source.Route), route.C)
			if err != nil {
				panic(err)
			}
			deleteSet[c.Id] = struct{}{}
		}
		if c.Source.Id == id {
			route, err := s.blocks[c.Target.Id].Block.GetRoute(core.RouteID(c.Target.Route))
			if err != nil {
				panic(err)
			}
			err = b.Block.Disconnect(core.RouteID(c.Source.Route), route.C)
			if err != nil {
				panic(err)
			}
			deleteSet[c.Id] = struct{}{}
		}
	}

	// delete the connections that involve this block
	for k, _ := range deleteSet {
		s.websocketBroadcast(Update{Action: DELETE, Type: CONNECTION, Data: s.connections[k]})
		delete(s.connections, k)
	}

	// stop and delete the block
	s.supervisor.Remove(b.Token)
	s.websocketBroadcast(Update{Action: DELETE, Type: BLOCK, Data: s.blocks[id]})
	delete(s.blocks, id)

	w.WriteHeader(http.StatusNoContent)

}
func (s *Server) BlockModify(w http.ResponseWriter, r *http.Request) {

}

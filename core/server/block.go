package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nikhan/go-fetch"
	"github.com/nytlabs/st-core/core"
	"github.com/thejerf/suture"
)

type BlockLedger struct {
	Name        string              `json:"name"`
	Type        string              `json:"type"`
	Id          int                 `json:"id"`
	Block       *core.Block         `json:"-"`
	Token       suture.ServiceToken `json:"-"`
	Parent      int                 `json:"parent"`
	Composition int                 `json:"composition,omitempty"`
	Routes      []core.Route        `json:"routes"`
	Outputs     []core.Output       `json:"outputs"`
}

type BlockUpdateRoute struct {
	Value struct {
		Fetch interface{} `json:"fetch,omitempty"`
		Json  interface{} `json:"json,omitempty"`
	} `json:"value"`
	Id    int `json:"id"`
	Route int `json:"route"`
}

type BlockUpdateName struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type BlockUpdateGroup struct {
	Group int `json:"group"`
	Id    int `json:"id"`
}

func (s *Server) ListBlocks() []BlockLedger {
	var blocks []BlockLedger
	for _, b := range s.blocks {
		blocks = append(blocks, *b)
	}
	return blocks
}

func (s *Server) BlockIndex(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()

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

	s.Lock()
	defer s.Unlock()

	blockSpec, ok := s.library[m.Type]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"spec not found"})
		return
	}

	m.Id = s.GetNextID()
	m.Block = core.NewBlock(blockSpec)
	m.Token = s.supervisor.Add(m.Block)
	m.Routes = m.Block.GetRoutes()
	m.Outputs = m.Block.GetOutputs()
	s.blocks[m.Id] = &m
	s.websocketBroadcast(Update{Action: CREATE, Type: BLOCK, Data: m})
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) BlockModifyRoute(w http.ResponseWriter, r *http.Request) {
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

	routes, ok := vars["index"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"no route index supplied"})
		return
	}

	route, err := strconv.Atoi(routes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"could not read request body"})
		return
	}

	s.Lock()
	defer s.Unlock()

	b, ok := s.blocks[id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"block not found"})
		return
	}

	var v map[string]interface{}
	err = json.Unmarshal(body, &v)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"could not unmarshal value"})
		return
	}

	var m interface{}
	_, isFetch := v["fetch"]
	if isFetch {
		queryString, ok := v["fetch"].(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			writeJSON(w, Error{"fetch is not string"})
			return
		}

		fo, err := fetch.Parse(queryString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJSON(w, Error{err.Error()})
			return
		}

		m = fo
	}

	_, isJSON := v["json"]
	if isJSON {
		m = v["json"]
	}

	if !isJSON && !isFetch {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"no value or query specified"})
		return
	}

	if isJSON || isFetch {
		err := b.Block.SetRoute(core.RouteID(route), m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJSON(w, Error{err.Error()})
			return
		}
	}

	rs, _ := s.blocks[id].Block.GetRoute(core.RouteID(route))
	s.blocks[id].Routes[route] = rs

	update := BlockUpdateRoute{
		Id:    id,
		Route: route,
	}

	if isFetch {
		update.Value.Fetch = m.(*fetch.Query).String()
	} else {
		update.Value.Json = m
	}

	s.websocketBroadcast(Update{Action: UPDATE, Type: BLOCK, Data: update})
	w.WriteHeader(http.StatusNoContent)
}
func (s *Server) BlockModifyName(w http.ResponseWriter, r *http.Request) {
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"could not read request body"})
		return
	}

	s.Lock()
	defer s.Unlock()

	_, ok = s.blocks[id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"block not found"})
		return
	}

	var name string
	err = json.Unmarshal(body, &name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"could not unmarshal value"})
		return
	}

	s.blocks[id].Name = name

	update := BlockUpdateName{
		Id:   id,
		Name: name,
	}

	s.websocketBroadcast(Update{Action: UPDATE, Type: BLOCK, Data: update})
	w.WriteHeader(http.StatusNoContent)
}
func (s *Server) BlockModifyGroup(w http.ResponseWriter, r *http.Request) {
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
		writeJSON(w, Error{"block not found"})
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

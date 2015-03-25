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

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type ProtoBlock struct {
	Label    string   `json:"label"`
	Parent   int      `json:"parent"`
	Type     string   `json:"type"`
	Position Position `json:"position"`
}

type BlockLedger struct {
	Label       string        `json:"label"`
	Type        string        `json:"type"`
	Id          int           `json:"id"`
	Block       *core.Block   `json:"-"`
	Parent      *Group        `json:"-"`
	Composition int           `json:"composition,omitempty"`
	Inputs      []core.Input  `json:"inputs"`
	Outputs     []core.Output `json:"outputs"`
	Position    Position      `json:"position"`
}

func (bl *BlockLedger) GetID() int {
	return bl.Id
}

func (bl *BlockLedger) GetParent() *Group {
	return bl.Parent
}

func (bl *BlockLedger) SetParent(group *Group) {
	bl.Parent = group
}

func (s *Server) ListBlocks() []BlockLedger {
	blocks := []BlockLedger{}
	for _, b := range s.blocks {
		blocks = append(blocks, *b)
	}
	return blocks
}

func (s *Server) BlockIndexHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(s.ListBlocks()); err != nil {
		panic(err)
	}
}

func (s *Server) BlockHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromMux(mux.Vars(r))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, err)
		return
	}

	s.Lock()
	defer s.Unlock()

	b, ok := s.blocks[id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"could not find block"})
		return
	}

	w.WriteHeader(http.StatusOK)
	writeJSON(w, b)
	return
}

func (s *Server) BlockModifyPositionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromMux(mux.Vars(r))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"could not read request body"})
		return
	}

	var p Position
	err = json.Unmarshal(body, &p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"could not read JSON"})
		return
	}

	s.Lock()
	defer s.Unlock()

	b, ok := s.blocks[id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"could not find block"})
		return
	}

	b.Position = p

	s.websocketBroadcast(Update{Action: UPDATE, Type: BLOCK, Data: wsBlock{wsPosition{wsId{id}, p}}})
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) CreateBlock(p ProtoBlock) (*BlockLedger, error) {
	blockSpec, ok := s.library[p.Type]
	if !ok {
		return nil, errors.New("spec not found")
	}

	block := core.NewBlock(blockSpec)

	m := &BlockLedger{
		Label:    p.Label,
		Position: p.Position,
		Type:     p.Type,
		Block:    block,
		Id:       s.GetNextID(),
	}

	if _, ok := s.groups[p.Parent]; !ok {
		return nil, errors.New("invalid group, could not create block")
	}

	go block.Serve()
	m.Inputs = block.GetInputs()
	m.Outputs = block.GetOutputs()
	s.blocks[m.Id] = m

	s.websocketBroadcast(Update{Action: CREATE, Type: BLOCK, Data: wsBlock{*m}})

	err := s.AddChildToGroup(p.Parent, m)
	if err != nil {
		return nil, err

	}

	return m, nil
}

// CreateBlockHandler responds to a POST request to instantiate a new block and add it to the Server.
func (s *Server) BlockCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"could not read request body"})
		return
	}

	var m ProtoBlock
	err = json.Unmarshal(body, &m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"no ID supplied"})
		return
	}

	s.Lock()
	defer s.Unlock()

	b, err := s.CreateBlock(m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	writeJSON(w, b)
}

func (s *Server) BlockModifyNameHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromMux(mux.Vars(r))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, err)
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

	_, ok := s.blocks[id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"block not found"})
		return
	}

	var label string
	err = json.Unmarshal(body, &label)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"could not unmarshal value"})
		return
	}

	s.blocks[id].Label = label

	s.websocketBroadcast(Update{Action: UPDATE, Type: BLOCK, Data: wsBlock{wsLabel{wsId{id}, label}}})
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) DeleteBlock(id int) error {
	b, ok := s.blocks[id]
	if !ok {
		return errors.New("block not found")
	}

	deleteSet := make(map[int]struct{})

	// build a set of connections that we may need to delete
	// we need to panic here because if any error is thrown we are in huge trouble
	// any panic indicates that our server connection ledger is no longer true
	for _, c := range s.connections {
		if c.Target.Id == id || c.Source.Id == id {
			deleteSet[c.Id] = struct{}{}
		}
	}

	// delete the connections that involve this block
	for k, _ := range deleteSet {
		s.DeleteConnection(k)
	}

	// remove from group
	s.DetachChild(b)

	// stop and delete the block
	b.Block.Stop()
	s.websocketBroadcast(Update{Action: DELETE, Type: BLOCK, Data: wsBlock{wsId{id}}})
	delete(s.blocks, id)
	return nil
}

func (s *Server) BlockDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromMux(mux.Vars(r))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, err)
		return
	}

	s.Lock()
	defer s.Unlock()

	err = s.DeleteBlock(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) BlockModifyRouteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := getIDFromMux(vars)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, err)
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

	var v *core.InputValue
	err = json.Unmarshal(body, &v)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"could not unmarshal value"})
		return
	}

	s.Lock()
	defer s.Unlock()

	err = s.ModifyBlockRoute(id, route, v)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) ModifyBlockRoute(id int, route int, v *core.InputValue) error {
	b, ok := s.blocks[id]
	if !ok {
		return errors.New("could not find block")
	}

	var value *core.InputValue

	if v.Exists() {
		value = v
	}

	err := b.Block.SetInput(core.RouteIndex(route), value)
	if err != nil {
		return err
	}

	s.blocks[id].Inputs[route].Value = value

	s.websocketBroadcast(Update{Action: UPDATE, Type: ROUTE, Data: wsRouteModify{ConnectionNode{id, route}, value}})
	return nil
}

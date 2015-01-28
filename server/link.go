package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type LinkLedger struct {
	Source int `json:"source"` // the source id
	Block  int `json:"block"`  // the block id
	Id     int `json:"id"`     // link id
}

type ProtoLink struct {
	Source int `json:"source"` // the source id
	Block  int `json:"block"`  // the block id
}

func (s *Server) CreateLink(l ProtoLink) (*LinkLedger, error) {
	b, ok := s.blocks[l.Block]
	if !ok {
		return nil, errors.New("could not find block")
	}

	sl, ok := s.sources[l.Source]
	if !ok {
		return nil, errors.New("could not find source")
	}

	link := &LinkLedger{
		Source: l.Source,
		Block:  l.Block,
		Id:     s.GetNextID(),
	}

	if b.GetParent() != sl.GetParent() {
		return nil, errors.New("block and source must be in the same group, cannot link")
	}

	err := b.Block.SetSource(sl.Source)
	if err != nil {
		return nil, err
	}

	s.links[link.Id] = link

	s.websocketBroadcast(Update{Action: CREATE, Type: LINK, Data: link})

	return link, nil
}

func (s *Server) DeleteLink(id int) error {
	link, ok := s.links[id]
	if !ok {
		return errors.New("could not find link")
	}

	s.blocks[link.Block].Block.SetSource(nil)
	delete(s.links, id)

	s.websocketBroadcast(Update{Action: CREATE, Type: CONNECTION, Data: link})
	return nil
}

func (s *Server) ListLinks() []LinkLedger {
	links := []LinkLedger{}
	for _, l := range s.links {
		links = append(links, *l)
	}
	return links
}

func (s *Server) LinkIndexHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()

	c := s.ListLinks()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(c); err != nil {
		panic(err)
	}
}

func (s *Server) LinkCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	var newLink ProtoLink
	json.Unmarshal(body, &newLink)

	s.Lock()
	defer s.Unlock()

	nl, err := s.CreateLink(newLink)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	writeJSON(w, nl)
}

func (s *Server) LinkDeleteHandler(w http.ResponseWriter, r *http.Request) {
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

	err = s.DeleteLink(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

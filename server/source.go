package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nytlabs/st-core/core"
	"github.com/thejerf/suture"
)

type SourceLedger struct {
	Label      string              `json:"label"`
	Type       string              `json:"type"`
	Id         int                 `json:"id"`
	Source     core.Source         `json:"-"`
	Parent     *Group              `json:"-"`
	Token      suture.ServiceToken `json:"-"`
	Position   Position            `json:"position"`
	Parameters []map[string]string `json:"params"`
}

type ProtoSource struct {
	Label      string            `json:"label"`
	Type       string            `json:"type"`
	Position   Position          `json:"position"`
	Parent     int               `json:"parent"`
	Parameters map[string]string `json:"params"`
}

func (sl *SourceLedger) GetID() int {
	return sl.Id
}

func (sl *SourceLedger) GetParent() *Group {
	return sl.Parent
}

func (sl *SourceLedger) SetParent(group *Group) {
	sl.Parent = group
}

func (s *Server) ListSources() []SourceLedger {
	sources := []SourceLedger{}
	for _, source := range s.sources {
		sources = append(sources, *source)
	}
	return sources
}

func (s *Server) SourceIndexHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(s.ListSources()); err != nil {
		panic(err)
	}
}

func (s *Server) SourceHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromMux(mux.Vars(r))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, err)
		return
	}
	s.Lock()
	defer s.Unlock()
	source, ok := s.sources[id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"could not find source"})
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJSON(w, source)
}

func (s *Server) CreateSource(p ProtoSource) (*SourceLedger, error) {
	f, ok := s.sourceLibrary[p.Type]
	if !ok {
		return nil, errors.New("source type " + p.Type + " does not exist")
	}

	source := f.New()

	sl := &SourceLedger{
		Label:      p.Label,
		Position:   p.Position,
		Source:     source,
		Type:       p.Type,
		Id:         s.GetNextID(),
		Parameters: make([]map[string]string, 0), // this will get overwritten if we have parameters
	}

	if i, ok := source.(core.Interface); ok {
		go i.Serve()
	}

	s.sources[sl.Id] = sl
	s.websocketBroadcast(Update{Action: CREATE, Type: SOURCE, Data: wsSource{*sl}})

	err := s.AddChildToGroup(p.Parent, sl)
	if err != nil {
		return nil, err

	}

	return sl, nil
}

func (s *Server) DeleteSource(id int) error {
	source, ok := s.sources[id]
	if !ok {
		return errors.New("could not find source")
	}

	for _, l := range s.links {
		if l.Source.Id == id {
			err := s.DeleteLink(l.Id)
			if err != nil {
				return err
			}
		}
	}

	if si, ok := source.Source.(core.Interface); ok {
		si.Stop()
	}

	s.DetachChild(source)

	s.websocketBroadcast(Update{Action: DELETE, Type: SOURCE, Data: wsSource{wsId{id}}})
	delete(s.sources, source.Id)
	return nil
}

func (s *Server) SourceCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"could not read request body"})
		return
	}

	var m ProtoSource
	err = json.Unmarshal(body, &m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"no ID supplied"})
		return
	}

	s.Lock()
	defer s.Unlock()

	b, err := s.CreateSource(m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	writeJSON(w, b)
}

func (s *Server) SourceDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromMux(mux.Vars(r))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, err)
		return
	}

	s.Lock()
	defer s.Unlock()

	err = s.DeleteSource(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func (s *Server) SourceGetValueHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromMux(mux.Vars(r))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, err)
		return
	}

	s.Lock()
	defer s.Unlock()

	val, err := s.GetSourceValue(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(val)
}

func (s *Server) SourceModifyPositionHandler(w http.ResponseWriter, r *http.Request) {
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

	b, ok := s.sources[id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"could not find block"})
		return
	}

	b.Position = p

	s.websocketBroadcast(Update{Action: UPDATE, Type: SOURCE, Data: wsSource{wsPosition{wsId{id}, p}}})
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) SourceModifyNameHandler(w http.ResponseWriter, r *http.Request) {
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

	_, ok := s.sources[id]
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

	s.sources[id].Label = label

	s.websocketBroadcast(Update{Action: UPDATE, Type: SOURCE, Data: wsSource{wsLabel{wsId{id}, label}}})
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) SourceSetValueHandler(w http.ResponseWriter, r *http.Request) {
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

	err = s.SetSourceValue(id, body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) GetSourceValue(id int) ([]byte, error) {
	source, ok := s.sources[id]
	if !ok {
		return nil, errors.New("source does not exist")
	}

	store, ok := source.Source.(core.Store)
	if !ok {
		return nil, errors.New("can only get values from stores")
	}

	store.Lock()
	defer store.Unlock()
	out, err := json.Marshal(store.Get())
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Server) SetSourceValue(id int, body []byte) error {
	source, ok := s.sources[id]
	if !ok {
		return errors.New("source does not exist")
	}

	store, ok := source.Source.(core.Store)
	if !ok {
		return errors.New("can only get values from stores")
	}

	var m interface{}
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}

	store.Lock()
	defer store.Unlock()
	err = store.Set(m)
	if err != nil {
		return err
	}

	return nil
}

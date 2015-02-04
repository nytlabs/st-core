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
	Parameters map[string]string   `json:"params"`
}

type ProtoSource struct {
	Label    string   `json:"label"`
	Type     string   `json:"type"`
	Position Position `json:"position"`
	Parent   int      `json:"group"`
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
		Label:    p.Label,
		Position: p.Position,
		Source:   source,
		Type:     p.Type,
		Id:       s.GetNextID(),
	}

	// Describe() is not thread-safe it must be put ahead of supervior...
	sl.Parameters = source.Describe()
	sl.Token = s.supervisor.Add(source)
	s.sources[sl.Id] = sl
	s.websocketBroadcast(Update{Action: CREATE, Type: SOURCE, Data: sl})

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
		if l.Source == id {
			err := s.DeleteLink(l.Id)
			if err != nil {
				return err
			}
		}
	}

	s.DetachChild(source)
	s.supervisor.Remove(source.Token)
	s.websocketBroadcast(Update{Action: DELETE, Type: SOURCE, Data: s.sources[id]})
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

func (s *Server) ModifySource(id int, m map[string]string) error {
	source, ok := s.sources[id]
	if !ok {
		return errors.New("no source found")
	}

	s.supervisor.Remove(source.Token)
	for k, _ := range source.Parameters {
		if v, ok := m[k]; ok {
			s.sources[id].Source.SetSourceParameter(k, v)
			source.Parameters[k] = v
			update := struct {
				Id    int    `json:"id"`
				Key   string `json:"param"`
				Value string `json:"value"`
			}{
				id, k, v,
			}
			s.websocketBroadcast(Update{Action: UPDATE, Type: SOURCE, Data: update})
		}
	}
	source.Token = s.supervisor.Add(source.Source)
	return nil
}

func (s *Server) SourceModifyHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"could not read request body"})
		return
	}

	var m map[string]string
	err = json.Unmarshal(body, &m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"no ID supplied"})
		return
	}

	id, err := getIDFromMux(mux.Vars(r))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, err)
		return
	}

	s.Lock()
	defer s.Unlock()

	err = s.ModifySource(id, m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

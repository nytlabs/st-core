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
	Parameters map[string]string   `json:"params,omitempty"`
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

	if i, ok := source.(core.Interface); ok {
		// Describe() is not thread-safe it must be put ahead of supervior...
		sl.Parameters = i.Describe()
		sl.Token = s.supervisor.Add(i)
	}

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

	if _, ok := source.Source.(core.Interface); ok {
		s.supervisor.Remove(source.Token)
	}

	s.DetachChild(source)

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

	i, ok := source.Source.(core.Interface)
	if !ok {
		return errors.New("cannot modify store")
	}

	s.supervisor.Remove(source.Token)
	for k, _ := range source.Parameters {
		if v, ok := m[k]; ok {
			i.SetSourceParameter(k, v)
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
	source.Token = s.supervisor.Add(i)
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

	update := struct {
		Position
		Id int
	}{
		p,
		id,
	}

	s.websocketBroadcast(Update{Action: UPDATE, Type: SOURCE, Data: update})
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

	update := struct {
		Id    int    `json:"id"`
		Label string `json:"label"`
	}{
		id, label,
	}

	s.websocketBroadcast(Update{Action: UPDATE, Type: SOURCE, Data: update})
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

package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Node interface {
	GetID() int
	GetParent() *Group
	SetParent(*Group)
}

type Group struct {
	Id       int    `json:"id"`
	Label    string `json:"label"`
	Children []int  `json:"children"`
	Parent   *Group `json:"-"`
}

func (g *Group) GetID() int {
	return g.Id
}

func (g *Group) GetParent() *Group {
	return g.Parent
}

func (g *Group) SetParent(group *Group) {
	g.Parent = group
}

func (s *Server) ListGroups() []Group {
	groups := []Group{}
	for _, g := range s.groups {
		groups = append(groups, *g)
	}
	return groups
}

func (s *Server) GroupIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(s.ListGroups()); err != nil {
		panic(err)
	}
}

func (s *Server) DetachChild(g Node) error {
	parent := g.GetParent()
	if parent == nil {
		return errors.New("no parent to detach from")
	}

	id := g.GetID()

	child := -1
	for i, v := range parent.Children {
		if v == id {
			child = i
		}
	}

	if child == -1 {
		return errors.New("could not remove child from group: child does not exist")
	}

	parent.Children = append(parent.Children[:child], parent.Children[child+1:]...)

	update := struct {
		Id    int `json:"id"`
		Child int `json:"child"`
	}{
		parent.GetID(), g.GetID(),
	}

	s.websocketBroadcast(Update{Action: DELETE, Type: GROUP_CHILD, Data: update})
	return nil
}

func (s *Server) AddChildToGroup(id int, n Node) error {
	newParent, ok := s.groups[id]
	if !ok {
		return errors.New("group not found")
	}

	nid := n.GetID()
	for _, v := range newParent.Children {
		if v == nid {
			return errors.New("node already child of this group")
		}
	}

	newParent.Children = append(newParent.Children, nid)
	if n.GetParent() != nil {
		err := s.DetachChild(n)
		if err != nil {
			return err
		}
	}

	n.SetParent(newParent)

	update := struct {
		Id    int `json:"id"`
		Child int `json:"child"`
	}{
		id, nid,
	}

	s.websocketBroadcast(Update{Action: UPDATE, Type: GROUP_CHILD, Data: update})
	return nil
}

// CreateGroupHandler responds to a POST request to instantiate a new group and add it to the Server.
// Moves all of the specified children out of the parent's group and into the new group.
func (s *Server) GroupCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"could not read request body"})
		return
	}

	var g struct {
		Group    int    `json:"group"`
		Children []int  `json:"children"`
		Label    string `json:"label"`
	}

	err = json.Unmarshal(body, &g)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"could not read JSON"})
		return
	}

	s.Lock()
	defer s.Unlock()

	newGroup := &Group{
		Children: g.Children,
		Label:    g.Label,
		Id:       s.GetNextID(),
	}

	if newGroup.Children == nil {
		newGroup.Children = []int{}
	}

	for _, c := range newGroup.Children {
		_, okb := s.blocks[c]
		_, okg := s.groups[c]
		if !okb && !okg {
			w.WriteHeader(http.StatusBadRequest)
			writeJSON(w, Error{"could not create group: invalid children"})
			return
		}
	}

	s.groups[newGroup.Id] = newGroup
	s.websocketBroadcast(Update{Action: CREATE, Type: GROUP, Data: newGroup})

	err = s.AddChildToGroup(g.Group, newGroup)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	for _, c := range newGroup.Children {
		if cb, ok := s.blocks[c]; ok {
			err = s.AddChildToGroup(newGroup.Id, cb)
		}
		if cg, ok := s.groups[c]; ok {
			err = s.AddChildToGroup(newGroup.Id, cg)
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJSON(w, Error{err.Error()})
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) DeleteGroup(id int) error {
	group, ok := s.groups[id]
	if !ok {
		return errors.New("could not find group to delete")
	}

	for _, c := range group.Children {
		if _, ok := s.blocks[c]; ok {
			err := s.DeleteBlock(c)
			if err != nil {
				return err
			}
		} else if _, ok := s.groups[c]; ok {
			err := s.DeleteGroup(c)
			if err != nil {
				return err
			}
		}
	}

	update := struct {
		Id int `json:"id"`
	}{
		id,
	}
	s.DetachChild(group)
	delete(s.groups, id)
	s.websocketBroadcast(Update{Action: DELETE, Type: GROUP, Data: update})
	return nil
}

func (s *Server) GroupDeleteHandler(w http.ResponseWriter, r *http.Request) {
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

	err = s.DeleteGroup(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func (s *Server) GroupHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) GroupExportHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) GroupImportHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) GroupModifyLabelHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) GroupModifyAllChildrenHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) GroupModifyChildHandler(w http.ResponseWriter, r *http.Request) {
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

	childs, ok := vars["node_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"no ID supplied"})
		return
	}

	child, err := strconv.Atoi(childs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	if id == child {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"cannot add group as member of itself"})
		return
	}

	s.Lock()
	defer s.Unlock()

	var n Node

	if _, ok := s.groups[id]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"could not find id"})
		return
	}

	if b, ok := s.blocks[child]; ok {
		n = b
	}
	if g, ok := s.groups[child]; ok {
		n = g
	}

	if n == nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"could not find id"})
		return
	}

	err = s.AddChildToGroup(id, n)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) GroupPositionHandler(w http.ResponseWriter, r *http.Request) {
}

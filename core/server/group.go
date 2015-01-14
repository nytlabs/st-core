package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
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
	oldParent := n.GetParent()

	// if this node had a previous parent assigned
	if oldParent != nil {
		child := -1
		for i, v := range oldParent.Children {
			if v == nid {
				child = i
			}
		}

		if child == -1 {
			return errors.New("could not remove child from group: child does not exist")
		}

		oldParent.Children = append(oldParent.Children[:child], oldParent.Children[child+1:]...)
	}

	n.SetParent(newParent)
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

	s.groups[newGroup.Id] = newGroup
	err = s.AddChildToGroup(g.Group, newGroup)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
	}

	s.websocketBroadcast(Update{Action: CREATE, Type: GROUP, Data: newGroup})
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) GroupDeleteHandler(w http.ResponseWriter, r *http.Request) {
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
}
func (s *Server) GroupPositionHandler(w http.ResponseWriter, r *http.Request) {
}

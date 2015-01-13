package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Group struct {
	Id       int    `json:"id"`
	Label    string `json:"label"`
	Children []int  `json:"children"`
	Parent   int    `json:"group"`
}

func (g *Group) GetID()

func (s *Server) ListGroups() []Group {
	group := []Group{}
	for _, g := range s.groups {
		groups = append(groups, *g)
	}
	return groups
}

func (s *Server) GroupIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(s.GroupList()); err != nil {
		panic(err)
	}
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

	g = &Group{
		Children: g.Children,
		Label:    g.Label,
		Id:       s.GetNextId(),
	}

	s.AddChildToGroup(g.Group, g.Id)
	s.Groups[g.Id] = g

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

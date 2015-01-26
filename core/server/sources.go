package server

import (
	"net/http"

	"github.com/nytlabs/st-core/core"
	"github.com/thejerf/suture"
)

type SourceLedger struct {
	Label    string              `json:"label"`
	Type     string              `json:"type"`
	Id       int                 `json:"id"`
	Block    *core.Source        `json:"-"`
	Parent   *Group              `json:"-"`
	Token    suture.ServiceToken `json:"-"`
	Position Position            `json:"position"`
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

func (s *Server) SourceIndexHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) SourceHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) SourceCreateHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) SourceModifyHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) SourceDeleteHandler(w http.ResponseWriter, r *http.Request) {
}

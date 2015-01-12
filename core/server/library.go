package server

import (
	"encoding/json"
	"net/http"
)

type LibraryEntry struct {
	Name   string `json:"name"`
	Source int    `json:"source"`
	// type if we need that later
}

func (s *Server) LibraryHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()

	l := []LibraryEntry{}

	for _, v := range s.library {
		l = append(l, LibraryEntry{
			v.Name,
			int(v.Shared),
		})
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(l); err != nil {
		panic(err)
	}
}

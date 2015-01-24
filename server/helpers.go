package server

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		panic(err)
	}
	return
}

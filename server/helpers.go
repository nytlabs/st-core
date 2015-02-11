package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		panic(err)
	}
	return
}

func getIDFromMux(vars map[string]string) (id int, err error) {
	ids, ok := vars["id"]
	if !ok {
		return 0, errors.New("no ID supplied")
	}

	id, err = strconv.Atoi(ids)
	if err != nil {
		return 0, err
	}
	return id, nil
}

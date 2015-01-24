package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEndpoints(t *testing.T) {

	s := NewServer()
	r := s.NewRouter()
	server := httptest.NewServer(r)
	defer server.Close()

	res, err := http.Get(server.URL + "/library")
	if err != nil {
		t.Error(err)
	}
	libraryJSON, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	var library interface{}
	err = json.Unmarshal(libraryJSON, &library)
	if err != nil {
		t.Error(err)
	}

}

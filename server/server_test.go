package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEndpoints(t *testing.T) {

	s := NewServer()
	r := s.NewRouter()
	server := httptest.NewServer(r)
	defer server.Close()
	routes := GetRoutes(s)
	for _, route := range routes {
		if route.Method != "GET" {
			continue
		}
		endpoint := route.Pattern
		endpoint = strings.Replace(endpoint, "{id}", "1", 1)
		res, err := http.Get(server.URL + endpoint)
		if err != nil {
			t.Error(err)
		}
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Error(err)
		}
		var marsh interface{}
		err = json.Unmarshal(body, &marsh)
		if err != nil {
			log.Println(string(body))
			t.Error(errors.New("failed to unmarshal response from " + endpoint))
		}
	}

}

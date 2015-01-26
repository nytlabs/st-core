package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestEndpoints(t *testing.T) {

	s := NewServer()
	r := s.NewRouter()
	server := httptest.NewServer(r)
	defer server.Close()

	// a couple of closures to save time below
	get := func(endpoint string) {
		res, err := http.Get(server.URL + endpoint)
		if err != nil {
			t.Error(err)
		}
		if res.StatusCode == 204 {
			return
		}
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Error(err)
		}
		var marsh interface{}
		err = json.Unmarshal(body, &marsh)
		if err != nil {
			t.Error(errors.New("GET: failed to unmarshal response from " + endpoint + ". Status code was: " + strconv.Itoa(res.StatusCode)))
		}
		b, err := json.MarshalIndent(marsh, "", "  ")
		if err != nil {
			t.Error("failed to Marshal")
		}
		fmt.Println(string(b) + "\n")
	}

	post := func(endpoint, msg string) {
		res, err := http.Post(server.URL+endpoint, "application/json", bytes.NewBuffer([]byte(msg)))
		if err != nil {
			t.Error(err)
		}
		if res.StatusCode == 204 {
			return
		}
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Error(err)
		}
		var marsh interface{}
		err = json.Unmarshal(body, &marsh)
		if err != nil {
			t.Error(errors.New("POST: failed to unmarshal response from " + endpoint + ". Status code was: " + strconv.Itoa(res.StatusCode)))
		}
		b, err := json.MarshalIndent(marsh, "", "  ")
		if err != nil {
			t.Error("failed to Marshal")
		}
		fmt.Println(string(b) + "\n")
	}

	put := func(endpoint, msg string) {
		req, err := http.NewRequest("PUT", server.URL+endpoint, bytes.NewBuffer([]byte(msg)))
		if err != nil {
			t.Error(err)
		}
		c := &http.Client{}
		res, err := c.Do(req)
		if err != nil {
			t.Error(err)
		}
		if res.StatusCode == 204 {
			return
		}
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Error(err)
		}
		var marsh interface{}
		err = json.Unmarshal(body, &marsh)
		if err != nil {
			t.Error(errors.New("PUT: failed to unmarshal response from " + endpoint + ". Response was: " + string(body)))
		}
		b, err := json.MarshalIndent(marsh, "", "  ")
		if err != nil {
			t.Error("failed to Marshal")
		}
		fmt.Println(string(b) + "\n")

	}

	// set up a group
	post("/groups", `{"group":0}`)

	// set up a + block
	post("/blocks", `{"type":"+","group":1}`)

	// label group 1
	put("/groups/1/label", "The Best Group Ever")

	// get all the groups
	get("/groups")

	// get group 1
	get("/groups/1")

	// get all the blocks
	get("/blocks")

	// get the + block
	get("/blocks/2")

	// export the pattern
	get("/groups/0/export")

}

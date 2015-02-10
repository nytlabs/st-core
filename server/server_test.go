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

	"github.com/fatih/color"
)

var warn = color.New(color.FgYellow).Add(color.Bold).Println

func procTestResponse(res *http.Response, t *testing.T, expectedCode int) {
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(body) + "\n")
	if res.StatusCode != expectedCode {
		t.Error(res.Request.Method + ": " + res.Request.URL.Path + " returned " + strconv.Itoa(res.StatusCode) + ". Expected " + strconv.Itoa(expectedCode) + ".")
		return
	}
	if res.StatusCode == 204 {
		return
	}
	if res.StatusCode == 404 {
		return
	}
	var marsh interface{}
	err = json.Unmarshal(body, &marsh)
	if err != nil {
		t.Error(errors.New(res.Request.Method + ": failed to unmarshal response from " + res.Request.URL.Path + ". Status code was: " + strconv.Itoa(res.StatusCode)))
	}
	_, err = json.MarshalIndent(marsh, "", "  ")
	if err != nil {
		t.Error("failed to Marshal")
	}
}

func TestEndpoints(t *testing.T) {

	s := NewServer()
	r := s.NewRouter()
	server := httptest.NewServer(r)
	defer server.Close()

	// a couple of closures to save time below
	get := func(endpoint string, expectedCode int) {
		res, err := http.Get(server.URL + endpoint)
		if err != nil {
			t.Error(err)
			return
		}
		procTestResponse(res, t, expectedCode)
	}

	post := func(endpoint, msg string, expectedCode int) {
		res, err := http.Post(server.URL+endpoint, "application/json", bytes.NewBuffer([]byte(msg)))
		if err != nil {
			t.Error(err)
			return
		}
		procTestResponse(res, t, expectedCode)
	}

	put := func(endpoint, msg string, expectedCode int) {
		req, err := http.NewRequest("PUT", server.URL+endpoint, bytes.NewBuffer([]byte(msg)))
		if err != nil {
			t.Error(err)
			return
		}
		c := &http.Client{}
		res, err := c.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		procTestResponse(res, t, expectedCode)
	}

	del := func(endpoint string, expectedCode int) {
		req, err := http.NewRequest("DELETE", server.URL+endpoint, nil)
		if err != nil {
			t.Error(err)
			return
		}
		c := &http.Client{}
		res, err := c.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		procTestResponse(res, t, expectedCode)
	}

	// set up a group (1)
	post("/groups", `{"group":0}`, 200)

	// set up a + block (2)
	post("/blocks", `{"type":"+","group":1}`, 200)

	// label group 1
	put("/groups/1/label", `"The Best Group Ever"`, 204)

	// label the plus block
	put("/blocks/2/label", `"my bestest adder"`, 204)

	// move the plus block
	put("/blocks/2/position", `{"x":10,"y":10}`, 204)

	// get all the groups
	get("/groups", 200)

	// get group 1
	get("/groups/1", 200)

	// get all the blocks
	get("/blocks", 200)

	// get the + block
	get("/blocks/2", 200)

	// make a delay block (3)
	post("/blocks", `{"type":"delay", "group":1}`, 200)

	// set the delay's value
	put("/blocks/3/routes/1", `{"type":"const","value":"1s"}`, 204)

	// make a log block (4)
	post("/blocks", `{"type":"log", "group":1}`, 200)

	// connect the + block to the delay block (5)
	post("/connections", `{"source":{"id":2, "Route":0}, "target":{"id":3, "Route":0}}`, 200)

	// connect the delay block to the log block (6)
	post("/connections", `{"source":{"id":3, "Route":0}, "target":{"id":4, "Route":0}}`, 200)

	// set the value of the plus inputs
	put("/blocks/2/routes/0", `{"type":"const","value":1}`, 204)
	put("/blocks/2/routes/1", `{"type":"const","value":1}`, 204)

	// make a set block (7)
	post("/blocks", `{"type":"set", "group":1}`, 200)

	// disconnect the log block from the delay block
	del("/connections/6", 204)

	// connect the set block to the log block and delay block (8) (9)
	post("/connections", `{"source":{"id":7, "Route":0}, "target":{"id":4, "Route":0}}`, 200)
	post("/connections", `{"source":{"id":3, "Route":0}, "target":{"id":7, "Route":1}}`, 200)

	// list connections
	get("/connections", 200)
	// describe connection 8
	get("/connections/8", 200)

	// set the value of the set key
	put("/blocks/7/routes/0", `{"type":"const","value":"myResult"}`, 204)

	// set the path on the log block
	put("/blocks/4/routes/0", `{"type":"fetch","value":".myResult"}`, 204)

	// move log block to root group
	put("/groups/0/children/4", "", 204)
	// move + block to root group (we will generate some errors with this later)
	put("/groups/0/children/2", "", 204)

	// create a keyvalue source (10)
	post("/sources", `{"type":"key-value"}`, 200)

	// get the keyvalue source
	get("/sources/10", 200)

	// make a stream source (11)
	post("/sources", `{"type":"stream"}`, 200)

	// change a parameter in the stream
	put("/sources/11", `{"topic":"test"}`, 204)

	// get all the sources
	get("/sources", 200)

	// make a key value get block (12)
	post("/blocks", `{"type":"kvGet"}`, 200)

	// link the key value get block to the key value source (13)
	post("/links", `{"source":10,"block":12}`, 200)

	// list the links
	get("/links", 200)

	// this doesn't exist yet - TODO use case?
	// get the link
	// get("/links/13", 200)

	// delete the link
	del("/links/13", 204)

	// delete the keyvalue store
	del("/sources/10", 204)

	// export the pattern
	get("/groups/0/export", 200)

	// import a pattern
	pattern := `{"blocks":[{"label":"","type":"+","id":2,"inputs":[{"name":"addend","type":"fetch","value":"."},{"name":"addend","type":"fetch","value":"."}],"outputs":[{"name":"sum"}],"position":{"x":0,"y":0}},{"label":"","type":"delay","id":3,"inputs":[{"name":"passthrough","type":"fetch","value":"."},{"name":"duration","type":"const","value":"1s"}],"outputs":[{"name":"passthrough"}],"position":{"x":0,"y":0}}],"connections":[{"source":{"id":2,"route":0},"target":{"id":3,"route":0},"id":4}],"groups":[{"id":1,"label":"","children":[2,3],"position":{"x":0,"y":0}}]}`
	post("/groups/1/import", pattern, 204)

	// delete the log block
	del("/blocks/4", 204)

	// delete group 1
	del("/groups/1", 204)

	// get the blocks library
	get("/blocks/library", 200)

	// get the blocks library
	get("/sources/library", 200)

	// generate some errors
	del("/groups/1", 400)                                                                       // delete a group we've already deleted
	del("/groups/", 404)                                                                        // delete unspecified group
	del("/blocks/246", 400)                                                                     // delete an unknown block
	post("/groups/1/import", "{}", 400)                                                         // import empty
	post("/groups/1/import", "{bla}", 400)                                                      // import malformed
	get("/groups/6/export", 400)                                                                // export an unknown group
	post("/sources", `{"type":"GodHead"}`, 400)                                                 // create an unknown source
	put("/groups/8/children/4", "", 400)                                                        // modify an unknown group
	put("/groups/0/children/34", "", 400)                                                       // move an unknown block to group 0
	put("/blocks/2/routes/20", `{"type":"fetch","value":".myResult"}`, 400)                     // set an unknown route
	put("/blocks/240/routes/0", `{"type":"fetch","value":".myResult"}`, 400)                    // set an unknown block's route
	put("/blocks/2/routes/0", `{"type":"fetch","value":"invalid"}`, 400)                        // set the + block's route to an invalid path
	put("/blocks/2/routes/0", `{"type":"value","value":"bob"}`, 400)                            // set the + block's route to an invalid value
	put("/blocks/2/routes/0", `{bobo}`, 400)                                                    // set the + block's route using malformed json
	post("/groups", `{"group":10}`, 400)                                                        // create a group with an unknown parent
	post("/groups", `{"group"10}`, 400)                                                         // create a group with malformed JSON
	post("/blocks", `{"type":"invalid", "group":0}`, 400)                                       // create a block of invalid type
	post("/blocks", `{"type"lid", "group":1}`, 400)                                             // create a block with malformed json
	post("/blocks", `{"type":"latch", "group":10}`, 400)                                        // create a block witha group that doesn't exist
	post("/connections", `{"source":{"id":700, "Route":0}, "target":{"id":2, "Route":0}}`, 400) //connect unknown source
	post("/connections", `{"source":{"id":2, "Route":0}, "target":{"id":200, "Route":0}}`, 400) //connect unknown target
	post("/connections", `{"source":{"i:0}, "ta200, "Route":0}}`, 400)                          //connect with malformed json
	post("/connections", `{}`, 400)                                                             //connect with empty json
	post("/connections", "", 400)                                                               //connect with empty string
	del("/connections/289", 400)                                                                //delete unknown connection
	del("/connections/", 404)                                                                   //delete unspecified connection
	del("/connections/invalid", 400)                                                            //delete malformed connection
}

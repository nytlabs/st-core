package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/fatih/color"
)

var warn = color.New(color.FgYellow).Add(color.Bold).Println

func procTestResponse(res *http.Response, t *testing.T) {
	if res.StatusCode == 204 {
		return
	}
	if res.StatusCode == 404 {
		warn("WARNING! " + res.Request.Method + ": " + res.Request.URL.Path + " returned 404")
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
		t.Error(errors.New(res.Request.Method + ": failed to unmarshal response from " + res.Request.URL.Path + ". Status code was: " + strconv.Itoa(res.StatusCode)))
	}
	_, err = json.MarshalIndent(marsh, "", "  ")
	if err != nil {
		t.Error("failed to Marshal")
	}
	// fmt.Println(string(b) + "\n")
}

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
		procTestResponse(res, t)
	}

	post := func(endpoint, msg string) {
		res, err := http.Post(server.URL+endpoint, "application/json", bytes.NewBuffer([]byte(msg)))
		if err != nil {
			t.Error(err)
		}
		procTestResponse(res, t)
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
		procTestResponse(res, t)
	}

	del := func(endpoint string) {
		req, err := http.NewRequest("DELETE", server.URL+endpoint, nil)
		if err != nil {
			t.Error(err)
		}
		c := &http.Client{}
		res, err := c.Do(req)
		if err != nil {
			t.Error(err)
		}
		procTestResponse(res, t)
	}

	// set up a group (1)
	post("/groups", `{"group":0}`)

	// set up a + block (2)
	post("/blocks", `{"type":"+","group":1}`)

	// label group 1
	put("/groups/1/label", "The Best Group Ever")

	// label the plus block
	put("/blocks/2/label", "my bestest adder")

	// get all the groups
	get("/groups")

	// get group 1
	get("/groups/1")

	// get all the blocks
	get("/blocks")

	// get the + block
	get("/blocks/2")

	// make a delay block (3)
	post("/blocks", `{"type":"delay", "group":1}`)

	// set the delay's value
	put("/blocks/3/routes/1", `{"type":"const","value":"1s"}`)

	// make a log block (4)
	post("/blocks", `{"type":"log", "group":1}`)

	// connect the + block to the delay block (5)
	post("/connections", `{"source":{"id":2, "Route":0}, "target":{"id":3, "Route":0}}`)

	// connect the delay block to the log block (6)
	post("/connections", `{"source":{"id":3, "Route":0}, "target":{"id":4, "Route":0}}`)

	// set the value of the plus inputs
	put("/blocks/2/routes/0", `{"type":"const","value":1}`)
	put("/blocks/2/routes/1", `{"type":"const","value":1}`)

	// make a set block (7)
	post("/blocks", `{"type":"set", "group":1}`)

	// disconnect the log block from the delay block
	del("/connections/6")

	// connect the set block to the log block and delay block (8) (9)
	post("/connections", `{"source":{"id":7, "Route":0}, "target":{"id":4, "Route":0}}`)
	post("/connections", `{"source":{"id":3, "Route":0}, "target":{"id":7, "Route":1}}`)

	// list connections
	get("/connections")
	// describe connection 8
	get("/connections/8")

	// set the value of the set key
	put("/blocks/7/routes/0", `{"type":"const","value":"myResult"}`)

	// set the path on the log block
	put("/blocks/4/routes/0", `{"type":"fetch","value":".myResult"}`)

	// move log block to root group
	put("/groups/0/children/4", "")
	// move + block to root group (we will generate some errors with this later)
	put("/groups/0/children/2", "")

	// create a keyvalue source (10)
	post("/sources", `{"type":"KeyValue"}`)

	// get all the sources
	get("/sources")

	// get the keyvalue source
	get("/sources/10")

	// make a stream source (11)
	post("/sources", `{"type":"Stream"}`)

	// change a parameter in the stream
	put("/sources/11", `{"topic":"test"}`)

	// delete the keyvalue store
	del("/sources/10")

	// export the pattern
	get("/groups/0/export")

	// import a pattern
	pattern := `{"blocks":[{"label":"","type":"+","id":2,"inputs":[{"name":"addend","type":"fetch","value":"."},{"name":"addend","type":"fetch","value":"."}],"outputs":[{"name":"sum"}],"position":{"x":0,"y":0}},{"label":"","type":"delay","id":3,"inputs":[{"name":"passthrough","type":"fetch","value":"."},{"name":"duration","type":"const","value":"1s"}],"outputs":[{"name":"passthrough"}],"position":{"x":0,"y":0}}],"connections":[{"source":{"id":2,"route":0},"target":{"id":3,"route":0},"id":4}],"groups":[{"id":1,"label":"","children":[2,3],"position":{"x":0,"y":0}}]}`
	post("/groups/1/import", pattern)

	// delete the log block
	del("/blocks/4")

	// delete group 1
	del("/groups/1")

	// get the library
	get("/library")

	// generate some errors
	del("/groups/1")                                                                       // delete a group we've already deleted
	del("/groups/")                                                                        // delete unspecified group
	del("/blocks/246")                                                                     // delete an unknown block
	post("/groups/1/import", "{}")                                                         // import empty
	post("/groups/1/import", "{bla}")                                                      // import malformed
	get("/groups/6/export")                                                                // export an unknown group
	post("/sources", `{"type":"GodHead"}`)                                                 // create an unknown source
	put("/groups/8/children/4", "")                                                        // modify an unknown group
	put("/groups/0/children/34", "")                                                       // move an unknown block to group 0
	put("/blocks/2/routes/20", `{"type":"fetch","value":".myResult"}`)                     // set an unknown route
	put("/blocks/240/routes/0", `{"type":"fetch","value":".myResult"}`)                    // set an unknown block's route
	put("/blocks/2/routes/0", `{"type":"fetch","value":"invalid"}`)                        // set the + block's route to an invalid path
	put("/blocks/2/routes/0", `{"type":"value","value":"bob"}`)                            // set the + block's route to an invalid value
	put("/blocks/2/routes/0", `{bobo}`)                                                    // set the + block's route using malformed json
	post("/groups", `{"group":10}`)                                                        // create a group with an unknown parent
	post("/groups", `{"group"10}`)                                                         // create a group with malformed JSON
	post("/blocks", `{"type":"invalid", "group":0}`)                                       // create a block of invalid type
	post("/blocks", `{"type"lid", "group":1}`)                                             // create a block with malformed json
	post("/blocks", `{"type":"latch", "group":10}`)                                        // create a block witha group that doesn't exist
	post("/connections", `{"source":{"id":700, "Route":0}, "target":{"id":2, "Route":0}}`) //connect unknown source
	//TODO this one panics
	post("/connections", `{"source":{"id":2, "Route":0}, "target":{"id":200, "Route":0}}`) //connect unknown target
	post("/connections", `{"source":{"i:0}, "ta200, "Route":0}}`)                          //connect with malformed json
	post("/connections", `{}`)                                                             //connect with empty json
	post("/connections", "")                                                               //connect with empty string
	del("/connections/289")                                                                //delete unknown connection
	del("/connections/")                                                                   //delete unspecified connection
	del("/connections/invalid")                                                            //delete malformed connection
}

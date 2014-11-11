package core

import (
	"log"

	"github.com/nikhan/go-fetch"
)

type Map struct {
	*Block
}

func NewMap(name string) Map {
	b := NewBlock(name)
	b.AddInput("in")
	b.AddInput("mapping")
	b.AddOutput("out")
	return Map{b}
}

func parseKeys(m map[string]interface{}) (interface{}, error) {
	t := make(map[string]interface{})

	for k, e := range m {
		switch r := e.(type) {
		case map[string]interface{}:
			//recurse
			j, err := parseKeys(r)
			if err != nil {
				return nil, err
			}
			t[k] = j
		case string:
			// this is a go-fetch directive
			q, err := fetch.Parse(r)
			if err != nil {
				return nil, err
			}
			t[k] = q
		}
	}

	return t, nil
}

func evalMap(msg interface{}, m map[string]interface{}) (interface{}, error) {

	t := make(map[string]interface{})
	for k, e := range m {
		switch r := e.(type) {
		case map[string]interface{}:
			j, err := evalMap(msg, r)
			if err != nil {
				return nil, err
			}
			t[k] = j
		case *fetch.Query:
			value, err := fetch.Run(r, msg)
			if err != nil {
				return nil, err
			}
			t[k] = value
		}
	}
	return t, nil

}

func (b Map) Serve() {

	in := b.GetInput("in")
	mapping := b.GetInput("mapping")
	var mI interface{}
	var p map[string]interface{}

	for {
		select {
		case msg := <-in.Connection:
			out, err := evalMap(msg, p)
			if err != nil {
				log.Fatal(err)
			}
			if ok := b.Broadcast(out, "out"); !ok {
				return
			}

		case mI = <-mapping.Connection:
			m, ok := mI.(map[string]interface{})
			if !ok {
				log.Fatal("could not assert mapping to map")
			}
			p, err := parseKeys(m)
			if err != nil {
				log.Fatal("could not parse keys")
			}
		case <-b.QuitChan:
			return
		}

	}
}

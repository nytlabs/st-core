package stores

import (
	"log"

	"github.com/nikhan/go-fetch"
	"github.com/nytlabs/st-core/core"
)

// KeyValue implements a simple key value store. Access is maintained via KeyValueGet and KeyValueSet blocks.
type KeyValue struct {
	*Store
	data    map[interface{}]interface{}
	setChan chan setRequest
	getChan chan getRequest
}

// block to Set a value of a KeyValueStore
type KeyValueSet struct {
	*core.Block
	Store *KeyValue
}

// block to Get a value from a KeyValueStore by key
type KeyValueGet struct {
	*core.Block
	Store *KeyValue
}

type setRequest struct {
	key   interface{}
	value interface{}
}

type getRequest struct {
	key      interface{}
	respChan chan interface{}
}

// returns a new KeyValue store
func NewKeyValue(name string) KeyValue {
	s := NewStore(name)
	return KeyValue{
		Store:   s,
		data:    make(map[interface{}]interface{}),
		setChan: make(chan setRequest),
		getChan: make(chan getRequest),
	}
}

func (s KeyValue) Serve() {
	for {
		select {
		case keyvalue := <-s.setChan:
			s.data[keyvalue.key] = keyvalue.value
		case req := <-s.getChan:
			req.respChan <- s.data[req.key]
		}
	}
}

// returns a new KeyValue Get block
func NewKeyValueGet(name string) KeyValueGet {
	b := core.NewBlock(name)
	b.AddInput("in")
	b.AddOutput("out")
	return KeyValueGet{
		Block: b,
	}
}

// connects a KeyValue Store to a KeyValue Get block
func (b KeyValueGet) ConnectStore(s *KeyValue) {
	b.Lock()
	b.Store = s
	b.Unlock()
}

func (b KeyValueGet) Serve() {

	in := b.GetInput("in")

	var v interface{}

	for {
		select {
		case m := <-in.Connection:
			k, err := fetch.Run(in.Path, m)
			if err != nil {
				log.Fatal(err)
			}
			r := getRequest{
				key:      k,
				respChan: make(chan interface{}),
			}
			select {
			case b.Store.getChan <- r:
			case <-b.QuitChan:
				return
			}
			select {
			case v = <-r.respChan:
			case <-b.QuitChan:
				return
			}
			if ok := b.Broadcast(v, "out"); !ok {
				return
			}

		case <-b.QuitChan:
			return
		}
	}

}

// returns a new KeyValue Set block
func NewKeyValueSet(name string) KeyValueSet {
	b := core.NewBlock(name)
	b.AddInput("in")
	return KeyValueSet{
		Block: b,
	}
}

// connections a KeyValue Set block to a KeyValue store
func (b KeyValueSet) ConnectStore(s *KeyValue) {
	b.Lock()
	b.Store = s
	b.Unlock()
}

func (b KeyValueSet) Serve() {

	in := b.GetInput("in")

	for {
		select {
		case m := <-in.Connection:
			d, ok := m.(map[interface{}]interface{})
			if !ok {
				log.Fatal("inbound message must be a map")
			}
			for k, v := range d {
				r := setRequest{
					key:   k,
					value: v,
				}
				b.Store.setChan <- r
			}
		}
	}

}

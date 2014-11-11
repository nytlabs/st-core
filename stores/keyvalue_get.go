package stores

import (
	"log"

	"github.com/nikhan/go-fetch"
	"github.com/nytlabs/st-core/core"
)

// block to Get a value from a KeyValueStore by key
type KeyValueGet struct {
	*core.Block
	Store *KeyValue
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

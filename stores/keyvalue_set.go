package stores

import (
	"log"

	"github.com/nikhan/go-fetch"
	"github.com/nytlabs/st-core/core"
)

// block to Set a value of a KeyValueStore
type KeyValueSet struct {
	*core.Block
	Store *KeyValue
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
			dI, err := fetch.Run(in.Path, m)
			if err != nil {
				log.Fatal("could not extract message")
			}
			d, ok := dI.(map[interface{}]interface{})
			if !ok {
				log.Fatal("must pass a map to KeyValueSet")
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

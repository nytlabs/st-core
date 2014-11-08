package core

import "log"

type KeyValueSet struct {
	*Block
	Store *KeyValue
}

func NewKeyValueSet(name string) KeyValueSet {
	b := NewBlock(name)
	b.AddInput("in")
	return KeyValueSet{
		Block: b,
	}
}

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

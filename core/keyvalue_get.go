package core

import (
	"log"

	"github.com/nikhan/go-fetch"
)

type KeyValueGet struct {
	*Block
	Store *KeyValue
}

func NewKeyValueGet(name string) KeyValueGet {
	b := NewBlock(name)
	b.AddInput("in")
	b.AddOutput("out")
	return KeyValueGet{
		Block: b,
	}
}

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

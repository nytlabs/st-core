package core

import (
	"log"
	"testing"
)

func TestStreamSend(t *testing.T) {
	log.Println("testing stream send")
	lib := GetLibrary()
	stream := NewStream()
	to := NewBlock(lib["streamSend"])
	go to.Serve()
	to.SetStore(stream)
	toRoute, _ := to.GetRoute(0)
	toRoute.C <- "hello from test stream send!"
}

func TestStreamReceive(t *testing.T) {
	log.Println("testing stream receive")
	lib := GetLibrary()
	stream := NewStream()
	from := NewBlock(lib["streamReceive"])
	go from.Serve()
	from.SetStore(stream)
	out := make(chan Message)
	from.Connect(0, out)
	m := <-out
	msg, ok := m.(map[string]string)
	if !ok {
		t.Error("didn't receive expected message")
	}
	_, ok = msg["hello"]
	if !ok {
		t.Error("didn't receive expected message")
	}
}

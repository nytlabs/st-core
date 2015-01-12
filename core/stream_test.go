package core

import (
	"log"
	"testing"
)

func TestStreamSend(t *testing.T) {
	log.Println("testing stream send")
	lib := GetLibrary()

	stream := NewStream()
	s.SetSourceParameter("topic", "test")
	s.SetSourceParameter("channel", "test#ephemeral")
	go stream.Serve()

	to := NewBlock(lib["streamSend"])
	go to.Serve()
	to.SetStore(stream)
	log.Println("sending")
	toRoute, _ := to.GetRoute(0)
	toRoute.C <- "hello from test stream send!"
	log.Println("sent")
	stream.Stop()
}

func TestStreamReceive(t *testing.T) {
	log.Println("testing stream receive")
	lib := GetLibrary()
	stream := NewStream()
	go stream.Serve()
	from := NewBlock(lib["streamReceive"])
	go from.Serve()
	from.SetStore(stream)
	out := make(chan Message)
	from.Connect(0, out)
	log.Println("receive is waiting for message")
	m := <-out
	msg, ok := m.(map[string]string)
	if !ok {
		t.Error("didn't receive expected message")
	}
	_, ok = msg["hello"]
	if !ok {
		t.Error("didn't receive expected message")
	}
	stream.Stop()
}

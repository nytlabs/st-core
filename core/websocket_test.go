package core

import (
	"log"
	"reflect"
	"testing"
)

func TestWsClient(t *testing.T) {

	log.Println("testing websocket client")

	wsSource := NewWsClient()
	ws, ok := wsSource.(Interface)
	if !ok {
		t.Fatal("could not assert websocket to Interface")
	}
	go ws.Serve()
	if ws.GetType() != WSCLIENT {
		t.Fatal("websocket client returns wrong type")
	}

	library := GetLibrary()
	blocks := map[string]*Block{
		"sink":            NewBlock(library["sink"]),
		"wsClientConnect": NewBlock(library["wsClientConnect"]),
		"wsClientReceive": NewBlock(library["wsClientReceive"]),
		"wsClientSend":    NewBlock(library["wsClientSend"]),
	}

	for _, v := range blocks {
		go v.Serve()
		go DummyMonitor(v.Monitor)
	}

	// hook up the client blocks
	err := blocks["wsClientConnect"].SetSource(ws)
	if err != nil {
		t.Fatal(err)
	}
	err = blocks["wsClientReceive"].SetSource(ws)
	if err != nil {
		t.Fatal(err)
	}
	err = blocks["wsClientSend"].SetSource(ws)
	if err != nil {
		t.Fatal(err)
	}

	// send before connecting ws
	sendOut := make(chan Message)
	sendIn, err := blocks["wsClientSend"].GetInput(0)
	if err != nil {
		t.Fatal(err)
	}
	err = blocks["wsClientSend"].Connect(0, sendOut)
	sendIn.C <- "should fail"
	if reflect.TypeOf(<-sendOut) != reflect.TypeOf(NewError("")) {
		t.Fatal("send should fail on closed websocket")
	}

	// connect together the rest of the pattern

	// 4. sink the connect and send
	sinkIn, err := blocks["sink"].GetInput(0)
	if err != nil {
		t.Fatal(err)
	}
	err = blocks["wsClientSend"].Connect(0, sinkIn.C)
	if err != nil {
		t.Fatal(err)
	}
	// 5. set the value of the connect block
	origin := InputValue{"http://localhost"}
	err = blocks["wsClientConnect"].SetInput(1, &origin)
	if err != nil {
		t.Fatal(err)
	}
	// 6. get ready to receive
	out := make(chan Message)
	err = blocks["wsClientReceive"].Connect(0, out)
	if err != nil {
		t.Fatal(err)
	}
	// 7. set the whole thing running
	urlIn, err := blocks["wsClientConnect"].GetInput(0)
	if err != nil {
		t.Fatal(err)
	}
	urlIn.C <- "ws://echo.websocket.org/"
	// block till we get something from the connect block
	connected := make(chan Message)
	err = blocks["wsClientConnect"].Connect(0, connected)
	if err != nil {
		t.Fatal(err)
	}
	if true != <-connected {
		t.Fatal("expected true from websocket connect")
	}
	// 8. send to websocket echo server
	testMessage := "howdy from streamtools!"
	sendIn.C <- testMessage
	sent := <-sendOut
	e, ok := sent.(*stcoreError)
	if ok {
		t.Fatal("got error from send block:", e)
	}
	got := <-out
	if testMessage != got {
		t.Fatal("expected different from websocket receive")
	}
	// stop the wsClient
	ws.Stop()
	sendIn.C <- "should fail"
	sent = <-sendOut
	_, ok = sent.(*stcoreError)
	if !ok {
		t.Fatal("send should fail on closed websocket", got)
	}

}

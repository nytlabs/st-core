package core

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/nsqio/nsq/nsqd"
	"github.com/nsqio/nsq/nsqlookupd"
)

// we need a bit of magic first, to turn on some nsq stuff (thanks snakes)
func bootstrapNSQCluster(t *testing.T) (string, []*nsqd.NSQD, []*nsqlookupd.NSQLookupd) {

	nsqlookupdOpts := nsqlookupd.NewOptions()
	nsqlookupdOpts.TCPAddress = "127.0.0.1:0"
	nsqlookupdOpts.HTTPAddress = "127.0.0.1:0"
	nsqlookupdOpts.BroadcastAddress = "127.0.0.1"
	nsqlookupd1 := nsqlookupd.New(nsqlookupdOpts)
	go nsqlookupd1.Main()

	time.Sleep(100 * time.Millisecond)

	nsqdOpts := nsqd.NewOptions()
	nsqdOpts.TCPAddress = "127.0.0.1:0"
	nsqdOpts.HTTPAddress = "127.0.0.1:0"
	nsqdOpts.BroadcastAddress = "127.0.0.1"
	nsqdOpts.NSQLookupdTCPAddresses = []string{nsqlookupd1.RealTCPAddr().String()}
	tmpDir, err := ioutil.TempDir("", fmt.Sprintf("nsq-test-%d", time.Now().UnixNano()))
	if err != nil {
		panic(err)
	}
	nsqdOpts.DataPath = tmpDir
	nsqd1 := nsqd.New(nsqdOpts)
	go nsqd1.Main()

	time.Sleep(100 * time.Millisecond)

	return tmpDir, []*nsqd.NSQD{nsqd1}, []*nsqlookupd.NSQLookupd{nsqlookupd1}
}

type TestNSQ struct {
}

func TestNSQConsumer(t *testing.T) {

	log.Println("testing nsq consumer")

	dataPath, nsqds, nsqlookupds := bootstrapNSQCluster(t)
	defer os.RemoveAll(dataPath)
	defer nsqds[0].Exit()
	defer nsqlookupds[0].Exit()

	// post a message
	topicName := "test"
	buf := bytes.NewBuffer([]byte("test message"))
	nsqdurl := fmt.Sprintf("http://%s/put?topic=%s", nsqds[0].RealHTTPAddr(), topicName)
	_, err := http.Post(nsqdurl, "application/octet-stream", buf)
	if err != nil {
		t.Fatal(err)
	}

	nsqlookupdurl := fmt.Sprintf("http://%s", nsqlookupds[0].RealHTTPAddr())

	nsqSource := NewNSQConsumer()
	nsq, ok := nsqSource.(Interface)
	if !ok {
		t.Fatal("could not assert nsq consumer to Interface")
	}
	go nsq.Serve()
	if nsq.GetType() != NSQCONSUMER {
		t.Fatal("nsq consumer returns wrong type")
	}

	library := GetLibrary()
	blocks := map[string]*Block{
		"log":     NewBlock(library["log"]),
		"connect": NewBlock(library["NSQConsumerConnect"]),
		"recv":    NewBlock(library["NSQConsumerReceive"]),
	}

	for _, v := range blocks {
		go v.Serve()
		go DummyMonitor(v.Monitor)
	}

	// hook up the client blocks
	err = blocks["connect"].SetSource(nsq)
	if err != nil {
		t.Fatal(err)
	}
	err = blocks["recv"].SetSource(nsq)
	if err != nil {
		t.Fatal(err)
	}

	// sink connect
	sinkIn, err := blocks["log"].GetInput(0)
	if err != nil {
		t.Fatal(err)
	}
	err = blocks["connect"].Connect(0, sinkIn.C)
	if err != nil {
		t.Fatal(err)
	}

	// listent to the receive
	out := make(chan Message)
	err = blocks["recv"].Connect(0, out)
	if err != nil {
		t.Fatal(err)
	}

	// connect the NSQ, relying on the fact that we haven't sunk the connect to only connect once
	topic := InputValue{topicName}
	channel := InputValue{"testChannel"}
	maxInFlight := InputValue{1.0}
	err = blocks["connect"].SetInput(0, &topic)
	if err != nil {
		t.Fatal(err)
	}
	err = blocks["connect"].SetInput(1, &channel)
	if err != nil {
		t.Fatal(err)
	}
	err = blocks["connect"].SetInput(3, &maxInFlight)
	if err != nil {
		t.Fatal(err)
	}

	in, err := blocks["connect"].GetInput(2)
	if err != nil {
		t.Fatal(err)
	}
	// fire
	in.C <- nsqlookupdurl

	// block on the receive's out chan for a moment
	select {
	case m := <-out:
		if m != "test message" {
			t.Fatal("received incorrect message")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("receive from NSQ timed out")
	}

	// stop the NSQ
	nsq.Stop()

}

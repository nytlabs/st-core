package core

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestList(t *testing.T) {
	log.Println("testing list")
	l := NewList()
	if l.GetType() != LIST {
		t.Fatal("list returns inaccurate id")
	}

	library := GetLibrary()
	blocks := map[string]*Block{
		"listGet":    NewBlock(library["listGet"]),
		"listSet":    NewBlock(library["listSet"]),
		"listAppend": NewBlock(library["listAppend"]),
		"listPop":    NewBlock(library["listPop"]),
		"listShift":  NewBlock(library["listShift"]),
		"listDump":   NewBlock(library["listDump"]),
	}

	for _, v := range blocks {
		go DummyMonitor(v.Monitor)
	}

	out := make(chan Message)
	for name, b := range blocks {
		log.Println("testing", name)
		go b.Serve()
		err := b.SetSource(l)
		if err != nil {
			t.Fatal(err)
		}
		b.Connect(0, out)
	}

	// put "foo" onto the empty list
	elementChan, err := blocks["listAppend"].GetInput(0)
	if err != nil {
		t.Fatal(err)
	}
	elementChan.C <- "foo"
	if <-out != true {
		t.Error("block did not produce expected value")
		return
	}

	// set the zeroth element to "bar"
	indexChan, err := blocks["listSet"].GetInput(0)
	if err != nil {
		t.Fatal(err)
	}
	elementChan, err = blocks["listSet"].GetInput(1)
	if err != nil {
		t.Fatal(err)
	}
	indexChan.C <- 0.0
	elementChan.C <- "bar"
	if <-out != true {
		t.Error("block did not produce expected value")
		return
	}

	// get the zeroth element and makes sure it's "bar"
	indexChan, err = blocks["listGet"].GetInput(0)
	if err != nil {
		t.Fatal(err)
	}
	indexChan.C <- 0.0
	if <-out != "bar" {
		t.Error("block did not produce expected value")
		return
	}

	// put "foo" into the front of the list
	elementChan, err = blocks["listShift"].GetInput(0)
	if err != nil {
		t.Fatal(err)
	}
	elementChan.C <- "foo"
	if <-out != true {
		t.Error("block did not produce expected value")
		return
	}

	// pop the last element off the list and make sure it's "bar"
	triggerChan, err := blocks["listPop"].GetInput(0)
	if err != nil {
		t.Fatal(err)
	}
	triggerChan.C <- true
	if <-out != "bar" {
		t.Error("block did not produce expected value")
		return
	}

	// dump the list and make sure it's ["foo"]
	triggerChan, err = blocks["listDump"].GetInput(0)
	if err != nil {
		t.Fatal(err)
	}
	triggerChan.C <- true
	msg := <-out
	myList := msg.([]interface{})

	if len(myList) != 1 {
		t.Error("dumped list wrong length")
		return
	}
	if myList[0] != "foo" {
		t.Error("dumped list has unexpected element")
	}
}

func TestValuePrimitive(t *testing.T) {
	log.Println("testing value primitive")
	v := NewValue()
	if v.GetType() != VALUE_PRIMITIVE {
		t.Fatal("Value source has wrong type")
	}
	library := GetLibrary()
	blocks := map[string]*Block{
		"valueGet": NewBlock(library["valueGet"]),
		"valueSet": NewBlock(library["valueSet"]),
	}

	for _, v := range blocks {
		go DummyMonitor(v.Monitor)
	}

	out := make(chan Message)
	for name, b := range blocks {
		log.Println("testing", name)
		go b.Serve()
		err := b.SetSource(v)
		if err != nil {
			t.Fatal(err)
		}
		b.Connect(0, out)
	}

	// set the value to "foo"
	valueChan, err := blocks["valueSet"].GetInput(0)
	if err != nil {
		t.Fatal(err)
	}
	valueChan.C <- "foo"
	if <-out != true {
		t.Error("block did not produce expected value")
		return
	}

	// make sure the value was set to "foo"
	triggerChan, err := blocks["valueGet"].GetInput(0)
	if err != nil {
		t.Fatal(err)
	}
	triggerChan.C <- true
	if <-out != "foo" {
		t.Error("block did not produce expected value")
		return
	}
}

func TestPriorityQueue(t *testing.T) {
	log.Println("testing priority queue")
	pq := NewPriorityQueue()
	if pq.GetType() != PRIORITY {
		t.Fatal("Priority Queue returns wrong type")
	}
	library := GetLibrary()
	blocks := map[string]*Block{
		"pqPush": NewBlock(library["pqPush"]),
		"pqPop":  NewBlock(library["pqPop"]),
		"pqPeek": NewBlock(library["pqPeek"]),
		"pqLen":  NewBlock(library["pqLen"]),
	}

	for _, v := range blocks {
		go DummyMonitor(v.Monitor)
	}

	out := make(chan Message)
	for name, b := range blocks {
		log.Println("testing", name)
		go b.Serve()
		err := b.SetSource(pq)
		if err != nil {
			t.Fatal(err)
		}
		b.Connect(0, out)
	}

	// connect up a channel for priorities for pqPop and pqPeek
	priorityChan := make(chan Message)

	blocks["pqPeek"].Connect(1, priorityChan)
	blocks["pqPop"].Connect(1, priorityChan)

	// push a message with the current time as priority
	p := float64(time.Now().Unix())
	in, err := blocks["pqPush"].GetInput(0)
	if err != nil {
		t.Fatal(err)
	}
	priority, err := blocks["pqPush"].GetInput(1)
	if err != nil {
		t.Fatal(err)
	}
	in.C <- "foo"
	priority.C <- p
	msg := <-out
	if msg != true {
		t.Fatal("pqPush did not generate expected output")
	}

	// peek at the message
	trigger, err := blocks["pqPeek"].GetInput(0)
	if err != nil {
		t.Fatal(err)
	}
	trigger.C <- true
	if <-out != "foo" {
		t.Fatal("pqPeek did not generate expected output")
	}
	if <-priorityChan != p {
		t.Fatal("pqPeek did not recive the expected priority")
	}

	// check the length of the queue
	lenTrigger, err := blocks["pqLen"].GetInput(0)
	if err != nil {
		t.Fatal(err)
	}
	lenTrigger.C <- true

	if l := <-out; l != 1.0 {
		length, ok := l.(float64)
		if !ok {
			t.Fatal("length produced something that's not a float")
		}
		t.Fatal(fmt.Sprint("pqLen did not generate expected output. Expecting 1 got", length))
	}

	// pop the message
	trigger, err = blocks["pqPop"].GetInput(0)
	if err != nil {
		t.Fatal(err)
	}
	trigger.C <- true
	if <-out != "foo" {
		t.Fatal("pqPop did not generate expected output")
	}
	if <-priorityChan != p {
		t.Fatal("pqPop did not recive the expected priority")
	}

	// check the length of the queue again to make sure pop did its thing
	lenTrigger.C <- true
	if <-out != 0.0 {
		t.Fatal("pqPop did not remove the message")
	}

}

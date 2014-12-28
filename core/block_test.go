package core

import (
	"encoding/json"
	"log"
	"testing"
	"time"
)

func TestDelay(t *testing.T) {
	log.Println("testing delay")
	spec := Delay()
	in := MessageMap{
		0: "test",
		1: "1s",
	}
	ic := make(chan Interrupt)
	out := MessageMap{}
	expected := MessageMap{0: "test"}
	tolerance, _ := time.ParseDuration("10ms")
	timerDuration, _ := time.ParseDuration("1s")
	timer := time.AfterFunc(timerDuration+tolerance, func() {
		t.Error("delay took longer than specified duration +", tolerance)
	})
	interrupt := spec.Kernel(in, out, nil, nil, ic)
	timer.Stop()
	if out[0] != expected[0] {
		t.Error("delay didn't pass the correct message")
	}
	if interrupt != nil {
		t.Error("delay returns inappropriate interrupt")
	}
}

func TestSingleBlock(t *testing.T) {
	log.Println("testing single block")

	out := make(chan Message)
	set := NewBlock(GetLibrary()["set"])
	go set.Serve()

	set.Connect(0, out)

	set.Input(0).C <- "testing"
	set.Input(1).C <- "success"

	p, err := json.Marshal(<-out)
	if err != nil {
		t.Error("could not marshal output of set block")
	}

	expected, _ := json.Marshal(map[string]interface{}{
		"testing": "success",
	})

	if string(p) != string(expected) {
		t.Error("not expected value")
	}

	set.Stop()
}

func TestKeyValue(t *testing.T) {
	log.Println("testing key value store")

	output := make(chan Message)
	sink := make(chan Message)

	testValues := map[string]string{
		"apple":      "red",
		"orange":     "orange",
		"pineapple":  "prickly",
		"grapefruit": "orange",
		"banana":     "yellow",
		"strawberry": "red",
	}

	kvset := NewBlock(GetLibrary()["kvSet"])
	kv := NewKeyValue()
	go kvset.Serve()
	kvset.Store(kv)

	kvset.Connect(0, sink)
	kvdump := NewBlock(GetLibrary()["kvDump"])
	go kvdump.Serve()
	kvdump.Store(kv)

	kvdump.Connect(0, output)
	for k, v := range testValues {
		kvset.Input(0).C <- k
		kvset.Input(1).C <- v
		_ = <-sink
	}

	kvdump.Input(0).C <- "bang"
	dump := <-output

	for k, vd := range dump.(map[string]Message) {
		if v, ok := testValues[k]; ok {
			if v != vd.(string) {
				t.Error("values not equal in kv store?!")
			}
		} else {
			t.Error("incomplete map in kv store")
		}
	}

	for k, vd := range testValues {
		if v, ok := dump.(map[string]Message)[k]; ok {
			if v.(string) != vd {
				t.Error("values not equal in kv store?!")
			}
		} else {
			t.Error("incomplete map in kv store")
		}
	}

}

func TestFirst(t *testing.T) {
	log.Println("testing first")
	f := NewBlock(GetLibrary()["first"])
	go f.Serve()
	sink := make(chan Message)
	f.Connect(0, sink)

	expected := []interface{}{true, false, false, false, false}
	in := f.Input(0).C

	for i, v := range expected {
		in <- i
		if v != <-sink {
			t.Error("first did not produce expected results")
		}
	}
}

func TestNull(t *testing.T) {
	log.Println("testing null stream")
	null := NewBlock(GetLibrary()["identity"])
	go null.Serve()
	null.RouteValue(0, nil)
	out := make(chan Message)
	null.Connect(0, out)
	o, err := json.Marshal(<-out)
	if err != nil {
		t.Error("could not marshall null stream")
	}
	if string(o) != "null" {
		t.Error("null stream is not null!")
	}
}

func BenchmarkAddition(b *testing.B) {
	sink := make(chan Message)
	add := NewBlock(GetLibrary()["+"])
	go add.Serve()
	add.Connect(0, sink)
	addend1 := add.Input(0).C
	addend2 := add.Input(1).C

	b.ResetTimer()
	for i := 0; i < 100000; i++ {
		addend1 <- 1.0
		addend2 <- 2.0
		_ = <-sink
	}
}

func BenchmarkRandomMath(b *testing.B) {
	sink := make(chan Message)
	u1 := NewBlock(GetLibrary()["uniform"])
	u2 := NewBlock(GetLibrary()["uniform"])
	u3 := NewBlock(GetLibrary()["uniform"])
	add := NewBlock(GetLibrary()["+"])
	mul := NewBlock(GetLibrary()["×"])
	go u1.Serve()
	go u2.Serve()
	go u3.Serve()
	go add.Serve()
	go mul.Serve()

	u1.Connect(0, add.Input(0).C)
	u2.Connect(0, add.Input(1).C)
	add.Connect(0, mul.Input(0).C)
	u3.Connect(0, mul.Input(1).C)
	mul.Connect(0, sink)

	b.ResetTimer()
	for i := 0; i < 100000; i++ {
		_ = <-sink
	}
}

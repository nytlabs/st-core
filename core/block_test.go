package core

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestSingleBlock(t *testing.T) {
	log.Println("testing single block")

	out := make(chan Message)
	set := NewBlock(GetLibrary()["set"])
	go set.Serve()

	set.Connect(0, out)

	sr1, _ := set.GetRoute(0)
	sr1.C <- "testing"
	sr2, _ := set.GetRoute(1)
	sr2.C <- "success"

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
	kvset.SetStore(kv)

	kvset.Connect(0, sink)
	kvdump := NewBlock(GetLibrary()["kvDump"])
	go kvdump.Serve()
	kvdump.SetStore(kv)

	kvdump.Connect(0, output)

	kv1, _ := kvset.GetRoute(0)
	kv2, _ := kvset.GetRoute(1)

	for k, v := range testValues {
		kv1.C <- k
		kv2.C <- v
		_ = <-sink
	}

	kvd, _ := kvdump.GetRoute(0)
	kvd.C <- "bang"

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

func TestRouteRace(t *testing.T) {
	sink := make(chan Message)
	identity := NewBlock(GetLibrary()["identity"])
	go identity.Serve()
	identity.Connect(0, sink)
	f := map[string]interface{}{
		"lol": "lol",
	}
	identity.SetRoute(0, f)

	z := <-sink

	if fmt.Sprintf("%p", f) == fmt.Sprintf("%p", z) {
		t.Error("route value race")
	}
}

func TestFirst(t *testing.T) {
	log.Println("testing first")
	f := NewBlock(GetLibrary()["first"])
	go f.Serve()
	sink := make(chan Message)
	f.Connect(0, sink)

	expected := []interface{}{true, false, false, false, false}
	in, _ := f.GetRoute(0)

	for i, v := range expected {
		in.C <- i
		if v != <-sink {
			t.Error("first did not produce expected results")
		}
	}
}

func TestNull(t *testing.T) {
	log.Println("testing null stream")
	null := NewBlock(GetLibrary()["identity"])
	go null.Serve()
	null.SetRoute(0, nil)
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
	addend1, _ := add.GetRoute(0)
	addend2, _ := add.GetRoute(1)

	b.ResetTimer()
	for i := 0; i < 100000; i++ {
		addend1.C <- 1.0
		addend2.C <- 2.0
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

	a1, _ := add.GetRoute(0)
	a2, _ := add.GetRoute(1)
	m1, _ := mul.GetRoute(0)
	m2, _ := mul.GetRoute(1)

	u1.Connect(0, a1.C)
	u2.Connect(0, a2.C)
	add.Connect(0, m1.C)
	u3.Connect(0, m2.C)
	mul.Connect(0, sink)

	b.ResetTimer()
	for i := 0; i < 100000; i++ {
		_ = <-sink
	}
}

package core

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
	"time"
)

func DummyMonitor(mm chan MonitorMessage) {
	for m := range mm {
		_ = m
	}
}

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
	go DummyMonitor(set.Monitor)
	go set.Serve()

	set.Connect(0, out)

	sr1, _ := set.GetInput(0)
	sr1.C <- "testing"
	sr2, _ := set.GetInput(1)
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
	go DummyMonitor(kvset.Monitor)
	kv := NewKeyValue()
	go kvset.Serve()
	kvset.SetSource(kv)

	kvset.Connect(0, sink)
	kvdump := NewBlock(GetLibrary()["kvDump"])
	go DummyMonitor(kvdump.Monitor)
	go kvdump.Serve()
	kvdump.SetSource(kv)

	kvdump.Connect(0, output)

	kv1, _ := kvset.GetInput(0)
	kv2, _ := kvset.GetInput(1)

	for k, v := range testValues {
		kv1.C <- k
		kv2.C <- v
		_ = <-sink
	}

	kvd, _ := kvdump.GetInput(0)
	kvd.C <- "bang"

	dump := <-output

	for k, vd := range dump.(map[string]interface{}) {
		if v, ok := testValues[k]; ok {
			if v != vd.(string) {
				t.Error("values not equal in kv store?!")
			}
		} else {
			t.Error("incomplete map in kv store")
		}
	}

	for k, vd := range testValues {
		if v, ok := dump.(map[string]interface{})[k]; ok {
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
	go DummyMonitor(f.Monitor)
	go f.Serve()
	sink := make(chan Message)
	f.Connect(0, sink)

	expected := []interface{}{true, false, false, false, false}
	in, _ := f.GetInput(0)

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
	go DummyMonitor(null.Monitor)
	go null.Serve()
	null.SetInput(0, &InputValue{nil})
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
	go DummyMonitor(add.Monitor)
	go add.Serve()
	add.Connect(0, sink)
	addend1, _ := add.GetInput(0)
	addend2, _ := add.GetInput(1)

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
	go DummyMonitor(u1.Monitor)
	u2 := NewBlock(GetLibrary()["uniform"])
	go DummyMonitor(u2.Monitor)
	u3 := NewBlock(GetLibrary()["uniform"])
	go DummyMonitor(u3.Monitor)
	add := NewBlock(GetLibrary()["+"])
	go DummyMonitor(add.Monitor)
	mul := NewBlock(GetLibrary()["×"])
	go DummyMonitor(mul.Monitor)
	go u1.Serve()
	go u2.Serve()
	go u3.Serve()
	go add.Serve()
	go mul.Serve()

	a1, _ := add.GetInput(0)
	a2, _ := add.GetInput(1)
	m1, _ := mul.GetInput(0)
	m2, _ := mul.GetInput(1)

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

func TestHTTPRequest(t *testing.T) {
	log.Println("testing HTTPRequest")
	lib := GetLibrary()
	block := NewBlock(lib["HTTPRequest"])
	go DummyMonitor(block.Monitor)
	go block.Serve()
	headers := make(map[string]interface{})
	block.SetInput(1, &InputValue{headers})
	block.SetInput(2, &InputValue{"GET"})
	block.SetInput(3, &InputValue{""})
	urlRoute, _ := block.GetInput(0)
	out := make(chan Message)
	block.Connect(0, out)
	urlRoute.C <- "http://private-e92ba-stcoretest.apiary-mock.com/get"
	m := <-out
	if reflect.DeepEqual(m, `{"msg": "hello there!"}`) {
		t.Error("didn't get expected output from HTTPRequest GET")
	}
}

func TestParseJSON(t *testing.T) {
	log.Println("testing ParseJSON")
	lib := GetLibrary()
	spec, ok := lib["parseJSON"]
	if !ok {
		t.Fatal("couldn't find block")
	}
	block := NewBlock(spec)
	go DummyMonitor(block.Monitor)
	go block.Serve()
	testJsonGood := "{\"foo\":\"bar\", \"weight\":2.3, \"someArray\":[1,2,3]}"
	out := make(chan Message)
	block.Connect(0, out)
	in, _ := block.GetInput(0)
	in.C <- testJsonGood
	m := <-out
	msg, ok := m.(map[string]interface{})
	if !ok {
		t.Error("expected map from ParseJSON")
		return
	}
	foo, ok := msg["foo"]
	if !ok {
		t.Error("missing expected key")
	}
	_, ok = foo.(string)
	if !ok {
		t.Error("expected string")
	}
	// now check it fails nicely
	testJsonBad := "{\"foo\":bar, \"weight\":2.3, \"someArray\":[1,2,3]}"
	in.C <- testJsonBad
	m = <-out
	_, ok = m.(error)
	if !ok {
		t.Error("expected error")
	}
}

func TestMerge(t *testing.T) {
	log.Println("testing merge")
	lib := GetLibrary()
	block := NewBlock(lib["merge"])
	go DummyMonitor(block.Monitor)
	go block.Serve()
	out := make(chan Message)
	inroute1, _ := block.GetInput(0)
	inroute2, _ := block.GetInput(1)
	block.Connect(0, out)
	inmsg1 := map[string]interface{}{"a": 3, "b": true}
	inmsg2 := map[string]interface{}{"c": 3}
	inmsg3 := map[string]interface{}{"b": "cat"}
	inmsg4 := map[string]interface{}{"a": 3, "b": true, "c": map[string]interface{}{"foo": false, "bar": "baz"}}
	inmsg5 := map[string]interface{}{"a": 3, "b": true, "c": map[string]interface{}{"foo": false, "bob": "bat"}}
	inmsg6 := map[string]interface{}{"a": 3, "b": true, "c": map[string]interface{}{"foo": false, "bar": "car"}}

	testMerge := func(i1, i2, expected map[string]interface{}) {
		inroute1.C <- i1
		inroute2.C <- i2
		ex, _ := json.Marshal(expected)
		got, err := json.Marshal(<-out)
		if err != nil {
			t.Error(err)
		}
		if string(got) != string(ex) {
			log.Println(string(got))
			log.Println(string(ex))
			t.Error("merge did not return the expected map")
		}
	}

	// simple test
	expected := map[string]interface{}{"a": 3, "b": true, "c": 3}
	testMerge(inmsg1, inmsg2, expected)

	// overwrite
	expected = map[string]interface{}{"a": 3, "b": "cat"}
	testMerge(inmsg1, inmsg3, expected)

	//nested
	expected = map[string]interface{}{"a": 3, "b": true, "c": map[string]interface{}{"foo": false, "bar": "baz", "bob": "bat"}}
	testMerge(inmsg4, inmsg5, expected)

	//nested overerite
	expected = map[string]interface{}{"a": 3, "b": true, "c": map[string]interface{}{"foo": false, "bar": "car"}}
	testMerge(inmsg4, inmsg6, expected)
}

package core

import (
	"log"
	"testing"
	"time"
)

type blockTest struct {
	in       MessageMap
	expected MessageMap
}

func TestSimpleDyads(t *testing.T) {
	simpleBlockTests := map[string]blockTest{
		"+": blockTest{
			in:       MessageMap{0: 1.0, 1: 2.0},
			expected: MessageMap{0: 3.0},
		},
		"-": blockTest{
			in:       MessageMap{0: 1.0, 1: 2.0},
			expected: MessageMap{0: -1.0},
		},
		"×": blockTest{
			in:       MessageMap{0: 3.0, 1: 2.0},
			expected: MessageMap{0: 6.0},
		},
		"÷": blockTest{
			in:       MessageMap{0: 6.0, 1: 2.0},
			expected: MessageMap{0: 3.0},
		},
		"^": blockTest{
			in:       MessageMap{0: 2.0, 1: 3.0},
			expected: MessageMap{0: 8.0},
		},
		"mod": blockTest{
			in:       MessageMap{0: 3.0, 1: 2.0},
			expected: MessageMap{0: 1.0},
		},
		">": blockTest{
			in:       MessageMap{0: 3.0, 1: 2.0},
			expected: MessageMap{0: true},
		},
		"<": blockTest{
			in:       MessageMap{0: 3.0, 1: 2.0},
			expected: MessageMap{0: false},
		},
		"==": blockTest{
			in:       MessageMap{0: 3.0, 1: "hello"},
			expected: MessageMap{0: false},
		},
		"!=": blockTest{
			in:       MessageMap{0: 3.0, 1: "hello"},
			expected: MessageMap{0: true},
		},
	}
	library := GetLibrary()
	for blockType, test := range simpleBlockTests {
		block, ok := library[blockType]
		if !ok {
			log.Fatal("could not find", blockType, "in library")
		}
		log.Println("testing", blockType)
		ic := make(chan Interrupt)
		out := MessageMap{}
		interrupt := block.Kernel(test.in, out, nil, nil, ic)
		for k, v := range test.expected {
			r, ok := out[k]
			if !ok {
				t.Error(blockType, "does not generate expected MessageMap")
			}
			if v != r {
				t.Error(blockType, "gives wrong output")
			}
		}
		if interrupt != nil {
			t.Error(blockType, "returns inappropriate interrupt")
		}
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

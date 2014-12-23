package core

import (
	"log"
	"testing"
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
	}
	library := GetLibrary()
	for blockType, test := range simpleBlockTests {
		block, ok := library[blockType]
		if !ok {
			log.Fatal("could not find", blockType, "in library")
		}
		ic := make(chan Interrupt)
		out := MessageMap{}
		interrupt := block.Kernel(test.in, out, ic)
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

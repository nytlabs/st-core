package core

import (
	"log"
	"math"
	"testing"
)

type blockTest struct {
	in       MessageMap
	expected MessageMap
}

func TestSimpleBlocks(t *testing.T) {
	simpleBlockTests := map[string]blockTest{
		"+": blockTest{
			in:       MessageMap{0: 1.0, 1: 2.0},
			expected: MessageMap{0: 3.0},
		},
		"-": blockTest{
			in:       MessageMap{0: 1.0, 1: 2.0},
			expected: MessageMap{0: -1.0},
		},
		"*": blockTest{
			in:       MessageMap{0: 3.0, 1: 2.0},
			expected: MessageMap{0: 6.0},
		},
		"/": blockTest{
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
		"exp": blockTest{
			in:       MessageMap{0: 3.0},
			expected: MessageMap{0: math.Exp(3.0)},
		},
		"log10": blockTest{
			in:       MessageMap{0: 3.0},
			expected: MessageMap{0: math.Log10(3.0)},
		},
		"ln": blockTest{
			in:       MessageMap{0: 3.0},
			expected: MessageMap{0: math.Log(3.0)},
		},
		"sqrt": blockTest{
			in:       MessageMap{0: 4.0},
			expected: MessageMap{0: math.Sqrt(4.0)},
		},
		"sin": blockTest{
			in:       MessageMap{0: 3.0},
			expected: MessageMap{0: math.Sin(3.0)},
		},
		"cos": blockTest{
			in:       MessageMap{0: 3.0},
			expected: MessageMap{0: math.Cos(3.0)},
		},
		"tan": blockTest{
			in:       MessageMap{0: 3.0},
			expected: MessageMap{0: math.Tan(3.0)},
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

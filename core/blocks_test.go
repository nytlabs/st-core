package core

import "testing"

func TestPlus(t *testing.T) {

	plusSpec := Plus()

	in := MessageMap{
		0: 1.0,
		1: 2.0,
	}

	out := MessageMap{}

	ic := make(chan Interrupt)

	interrupt := plusSpec.Kernel(in, out, ic)

	result, ok := out[0]

	if !ok {
		t.Error("plus does not generate expected MessageMap")
	}

	if result != 3.0 {
		t.Error("plus gives wrong output")
	}

	if interrupt != nil {
		t.Error("plus returns inappropriate interrupt")
	}

}

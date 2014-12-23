package core

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type stcoreError struct {
	s string
}

func (e *stcoreError) Error() string {
	log.Println(e.s)
	return e.s
}

func NewError(s string) *stcoreError {
	return &stcoreError{
		s: s,
	}
}

// Library is the set of all core block Specs
func GetLibrary() map[string]Spec {
	return map[string]Spec{
		"delay": Delay(),
		"set":   Set(),
		"log":   Log(),
		"sink":  Sink(),
		"+":     Addition(),
		"-":     Subtraction(),
		"×":     Multiplication(),
		"÷":     Division(),
	}
}

// Delay emits the message on passthrough after the specified duration
func Delay() Spec {
	return Spec{
		Inputs:  []Pin{Pin{"passthrough"}, Pin{"duration"}},
		Outputs: []Pin{Pin{"passthrough"}},
		Kernel: func(in MessageMap, out MessageMap, i chan Interrupt) Interrupt {
			t, err := time.ParseDuration(in[1].(string))
			if err != nil {
				out[0] = err
				return nil
			}
			timer := time.NewTimer(t)
			select {
			case <-timer.C:
				out[0] = in[0]
				return nil
			case f := <-i:
				return f
			}
			return nil
		},
	}
}

// Set creates a new message with the specified key and value
func Set() Spec {
	return Spec{
		Inputs:  []Pin{Pin{"key"}, Pin{"value"}},
		Outputs: []Pin{Pin{"object"}},
		Kernel: func(in MessageMap, out MessageMap, i chan Interrupt) Interrupt {
			out[0] = map[string]interface{}{
				in[0].(string): in[1],
			}
			return nil
		},
	}
}

// Log writes the inbound message to stdout
// TODO where should this write exactly?
func Log() Spec {
	return Spec{
		Inputs:  []Pin{Pin{"log"}},
		Outputs: []Pin{},
		Kernel: func(in MessageMap, out MessageMap, i chan Interrupt) Interrupt {
			o, err := json.Marshal(in[0])
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(o))
			return nil
		},
	}
}

// Sink discards the inbound message
func Sink() Spec {
	return Spec{
		Inputs:  []Pin{Pin{"in"}},
		Outputs: []Pin{},
		Kernel: func(in, out MessageMap, i chan Interrupt) Interrupt {
			return nil
		},
	}
}

// Addition returns the sum of the addenda
func Addition() Spec {
	return Spec{
		Inputs:  []Pin{Pin{"addend"}, Pin{"addend"}},
		Outputs: []Pin{Pin{"sum"}},
		Kernel: func(in, out MessageMap, i chan Interrupt) Interrupt {
			a1, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("Addition requires floats")
				return nil
			}
			a2, ok := in[1].(float64)
			if !ok {
				out[0] = NewError("Addition requires floats")
				return nil
			}
			out[0] = a1 + a2
			return nil
		},
	}
}

// Subtraction returns the difference of the minuend - subtrahend
func Subtraction() Spec {
	return Spec{
		Inputs:  []Pin{Pin{"minuend"}, Pin{"subtrahend"}},
		Outputs: []Pin{Pin{"difference"}},
		Kernel: func(in, out MessageMap, i chan Interrupt) Interrupt {
			m, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("Subtraction requires floats")
				return nil
			}
			s, ok := in[1].(float64)
			if !ok {
				out[0] = NewError("Subtraction requires floats")
				return nil
			}
			out[0] = m - s
			return nil
		},
	}
}

// Multiplication returns the product of the multiplicanda
func Multiplication() Spec {
	return Spec{
		Inputs:  []Pin{Pin{"multiplicand"}, Pin{"multiplicand"}},
		Outputs: []Pin{Pin{"product"}},
		Kernel: func(in, out MessageMap, i chan Interrupt) Interrupt {
			m1, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("Multiplication requires floats")
				return nil
			}
			m2, ok := in[1].(float64)
			if !ok {
				out[0] = NewError("Multiplication requires floats")
				return nil
			}
			out[0] = m1 * m2
			return nil
		},
	}
}

// Division returns the quotient of the dividend / divisor
func Division() Spec {
	return Spec{
		Inputs:  []Pin{Pin{"dividend"}, Pin{"divisor"}},
		Outputs: []Pin{Pin{"quotient"}},
		Kernel: func(in, out MessageMap, i chan Interrupt) Interrupt {
			d1, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("Division requires floats")
				return nil
			}
			d2, ok := in[1].(float64)
			if !ok {
				out[0] = NewError("Division requires floats")
				return nil
			}
			out[0] = d1 / d2
			return nil
		},
	}
}

package core

import (
	"encoding/json"
	"fmt"
	"time"
)

// Library is the set of all core block Specs
func GetLibrary() map[string]Spec {
	return map[string]Spec{
		"plus":  Plus(),
		"delay": Delay(),
		"set":   Set(),
		"log":   Log(),
	}
}

// Sum sums together numbers on each addend
func Plus() Spec {
	return Spec{
		Inputs: []Pin{
			Pin{"addend"},
			Pin{"addend"},
		},
		Outputs: []Pin{
			Pin{"sum"},
		},
		Kernel: func(in MessageMap, out MessageMap, i chan Interrupt) Interrupt {
			out[0] = in[0].(float64) + in[1].(float64)
			return nil
		},
	}
}

// Delay emits the message on passthrough after the specified duration
func Delay() Spec {
	return Spec{
		Inputs: []Pin{
			Pin{
				"passthrough",
			},
			Pin{
				"duration",
			},
		},
		Outputs: []Pin{
			Pin{
				"passthrough",
			},
		},
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
		Inputs: []Pin{
			Pin{
				"key",
			},
			Pin{
				"value",
			},
		},
		Outputs: []Pin{
			Pin{
				"object",
			},
		},
		Kernel: func(in MessageMap, out MessageMap, i chan Interrupt) Interrupt {
			out[0] = map[string]interface{}{
				in[0].(string): in[1],
			}
			return nil
		},
	}
}

// Log writes the inbound message to stdout TODO where should this write exactly?
func Log() Spec {
	return Spec{
		Inputs: []Pin{
			Pin{
				"log",
			},
		},
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

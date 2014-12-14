package core

import "time"

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
		Kernel: func(in MessageMap, out MessageMap, i chan InterruptFunc) InterruptFunc {
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
		Kernel: func(in MessageMap, out MessageMap, i chan InterruptFunc) InterruptFunc {
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

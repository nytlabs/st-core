package core

// sum together numbers on each `addend`, returning result on `sum`
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

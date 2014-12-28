package core

import "math"

// exp, log10, ln, sqrt, sin, cos, tan

func Exp() Spec {
	return Spec{
		Inputs:  []Pin{Pin{"power"}},
		Outputs: []Pin{Pin{"product"}},
		Kernel: func(in, out, internal MessageMap, s Store, i chan Interrupt) Interrupt {
			p, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("need float")
				return nil
			}
			out[0] = math.Exp(p)
			return nil
		},
	}
}

func Log10() Spec {
	return Spec{
		Inputs:  []Pin{Pin{"in"}},
		Outputs: []Pin{Pin{"log10"}},
		Kernel: func(in, out, internal MessageMap, s Store, i chan Interrupt) Interrupt {
			p, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("need float")
				return nil
			}
			out[0] = math.Log10(p)
			return nil
		},
	}
}
func Ln() Spec {
	return Spec{
		Inputs:  []Pin{Pin{"in"}},
		Outputs: []Pin{Pin{"ln"}},
		Kernel: func(in, out, internal MessageMap, s Store, i chan Interrupt) Interrupt {
			p, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("need float")
				return nil
			}
			out[0] = math.Log(p)
			return nil
		},
	}
}

func Sqrt() Spec {
	return Spec{
		Inputs:  []Pin{Pin{"in"}},
		Outputs: []Pin{Pin{"squareRoot"}},
		Kernel: func(in, out, internal MessageMap, s Store, i chan Interrupt) Interrupt {
			p, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("need float")
				return nil
			}
			out[0] = math.Sqrt(p)
			return nil
		},
	}
}

func Sin() Spec {
	return Spec{
		Inputs:  []Pin{Pin{"in"}},
		Outputs: []Pin{Pin{"sin"}},
		Kernel: func(in, out, internal MessageMap, s Store, i chan Interrupt) Interrupt {
			p, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("need float")
				return nil
			}
			out[0] = math.Sin(p)
			return nil
		},
	}
}

func Cos() Spec {
	return Spec{
		Inputs:  []Pin{Pin{"in"}},
		Outputs: []Pin{Pin{"cos"}},
		Kernel: func(in, out, internal MessageMap, s Store, i chan Interrupt) Interrupt {
			p, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("need float")
				return nil
			}
			out[0] = math.Cos(p)
			return nil
		},
	}
}

func Tan() Spec {
	return Spec{
		Inputs:  []Pin{Pin{"in"}},
		Outputs: []Pin{Pin{"tan"}},
		Kernel: func(in, out, internal MessageMap, s Store, i chan Interrupt) Interrupt {
			p, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("need float")
				return nil
			}
			out[0] = math.Tan(p)
			return nil
		},
	}
}

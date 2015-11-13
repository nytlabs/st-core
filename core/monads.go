package core

import "math"

// exp, log10, ln, sqrt, sin, cos, tan, floor, ceil

func Exp() Spec {
	return Spec{
		Name:    "exp",
		Inputs:  []Pin{Pin{"power", NUMBER}},
		Outputs: []Pin{Pin{"product", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
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

func Floor() Spec {
	return Spec{
		Name:    "floor",
		Inputs:  []Pin{Pin{"in", NUMBER}},
		Outputs: []Pin{Pin{"out", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			p, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("need float")
				return nil
			}
			out[0] = math.Floor(p)
			return nil
		},
	}
}

func Ceil() Spec {
	return Spec{
		Name:    "ceil",
		Inputs:  []Pin{Pin{"in", NUMBER}},
		Outputs: []Pin{Pin{"out", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			p, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("need float")
				return nil
			}
			out[0] = math.Ceil(p)
			return nil
		},
	}
}

func Log10() Spec {
	return Spec{
		Name:    "log10",
		Inputs:  []Pin{Pin{"in", NUMBER}},
		Outputs: []Pin{Pin{"log10", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
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
		Name:    "ln",
		Inputs:  []Pin{Pin{"in", NUMBER}},
		Outputs: []Pin{Pin{"ln", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
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
		Name:    "sqrt",
		Inputs:  []Pin{Pin{"in", NUMBER}},
		Outputs: []Pin{Pin{"squareRoot", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
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
		Name:    "sin",
		Inputs:  []Pin{Pin{"in", NUMBER}},
		Outputs: []Pin{Pin{"sin", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
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
		Name:    "cos",
		Inputs:  []Pin{Pin{"in", NUMBER}},
		Outputs: []Pin{Pin{"cos", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
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
		Name:    "tan",
		Inputs:  []Pin{Pin{"in", NUMBER}},
		Outputs: []Pin{Pin{"tan", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
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

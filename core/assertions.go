package core

func IsError() Spec {
	return Spec{
		Name:    "isError",
		Inputs:  []Pin{Pin{"in"}},
		Outputs: []Pin{Pin{"out"}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			_, ok := in[0].(stcoreError)
			if !ok {
				out[0] = false
				return nil
			}
			out[0] = true
			return nil
		},
	}
}

func IsBoolean() Spec {
	return Spec{
		Name:    "isBoolean",
		Inputs:  []Pin{Pin{"in"}},
		Outputs: []Pin{Pin{"out"}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			_, ok := in[0].(bool)
			if !ok {
				out[0] = false
				return nil
			}
			out[0] = true
			return nil
		},
	}
}

func IsString() Spec {
	return Spec{
		Name:    "isString",
		Inputs:  []Pin{Pin{"in"}},
		Outputs: []Pin{Pin{"out"}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			_, ok := in[0].(string)
			if !ok {
				out[0] = false
				return nil
			}
			out[0] = true
			return nil
		},
	}
}

func IsNumber() Spec {
	return Spec{
		Name:    "isNumber",
		Inputs:  []Pin{Pin{"in"}},
		Outputs: []Pin{Pin{"out"}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			_, ok := in[0].(float64)
			if !ok {
				out[0] = false
				return nil
			}
			out[0] = true
			return nil
		},
	}
}

func IsArray() Spec {
	return Spec{
		Name:    "isArray",
		Inputs:  []Pin{Pin{"in"}},
		Outputs: []Pin{Pin{"out"}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			_, ok := in[0].([]interface{})
			if !ok {
				out[0] = false
				return nil
			}
			out[0] = true
			return nil
		},
	}
}

func IsObject() Spec {
	return Spec{
		Name:    "isObject",
		Inputs:  []Pin{Pin{"in"}},
		Outputs: []Pin{Pin{"out"}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			_, ok := in[0].(map[string]interface{})
			if !ok {
				out[0] = false
				return nil
			}
			out[0] = true
			return nil
		},
	}
}

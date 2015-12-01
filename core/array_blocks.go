package core

// Head emits the first element of the inbound array
func Head() Spec {
	return Spec{
		Name:    "head",
		Inputs:  []Pin{Pin{"in", ARRAY}},
		Outputs: []Pin{Pin{"head", ANY}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			arr, ok := in[0].([]interface{})
			if !ok {
				out[0] = NewError("head requires an array")
				return nil
			}
			if len(arr) == 0 {
				out[0] = NewError("zero length array passed to head")
				return nil
			}
			out[0] = arr[0]
			return nil
		},
	}
}

// Tail emits all the elements of an array except for the first
func Tail() Spec {
	return Spec{
		Name:    "tail",
		Inputs:  []Pin{Pin{"in", ARRAY}},
		Outputs: []Pin{Pin{"tail", ARRAY}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			arr, ok := in[0].([]interface{})
			if !ok {
				out[0] = NewError("tail requires an array")
				return nil
			}
			if len(arr) == 0 {
				out[0] = NewError("zero length array passed to tail")
				return nil
			}
			out[0] = arr[1:len(arr)]
			return nil
		},
	}
}

// Last returns the last element of an array
func Last() Spec {
	return Spec{
		Name:    "last",
		Inputs:  []Pin{Pin{"in", ARRAY}},
		Outputs: []Pin{Pin{"last", ANY}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			arr, ok := in[0].([]interface{})
			if !ok {
				out[0] = NewError("last requires an array")
				return nil
			}
			if len(arr) == 0 {
				out[0] = NewError("zero length array passed to last")
				return nil
			}
			out[0] = arr[len(arr)-1]
			return nil
		},
	}
}

// Init returns the all the elements of an array apart from the last one
func Init() Spec {
	return Spec{
		Name:    "init",
		Inputs:  []Pin{Pin{"in", ARRAY}},
		Outputs: []Pin{Pin{"init", ARRAY}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			arr, ok := in[0].([]interface{})
			if !ok {
				out[0] = NewError("init requires an array")
				return nil
			}
			if len(arr) == 0 {
				out[0] = NewError("zero length array passed to init")
				return nil
			}
			out[0] = arr[0 : len(arr)-1]
			return nil
		},
	}
}

// Append appends the supplied element to the supplied array
func Append() Spec {
	return Spec{
		Name:    "append",
		Inputs:  []Pin{Pin{"element", ANY}, Pin{"array", ARRAY}},
		Outputs: []Pin{Pin{"array", ARRAY}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			arr, ok := in[1].([]interface{})
			if !ok {
				out[0] = NewError("Append requires an array")
				return nil
			}
			out[0] = append(arr, in[0])
			return nil
		},
	}
}

func Len() Spec {
	return Spec{
		Name:    "len",
		Inputs:  []Pin{Pin{"in", ARRAY}},
		Outputs: []Pin{Pin{"out", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			arr, ok := in[0].([]interface{})
			if !ok {
				out[0] = NewError("len requires an array")
				return nil
			}
			out[0] = float64(len(arr))
			return nil
		},
	}
}

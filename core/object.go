package core

// Keys returns the top level keys of the supplied object
func Keys() Spec {
	return Spec{
		Name:    "keys",
		Inputs:  []Pin{Pin{"in", OBJECT}},
		Outputs: []Pin{Pin{"keys", ARRAY}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			obj, ok := in[0].(map[string]interface{})
			if !ok {
				out[0] = NewError("Keys requires an object")
			}
			keys := make([]interface{}, len(obj))
			j := 0
			for k, _ := range obj {
				keys[j] = k
				j += 1
			}
			out[0] = keys
			return nil
		},
	}
}

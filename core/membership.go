package core

import "strings"

// InArray returns true if the supplied element is in the supplied array
func InArray() Spec {
	return Spec{
		Name:    "inArray",
		Inputs:  []Pin{Pin{"element", ANY}, Pin{"array", ARRAY}},
		Outputs: []Pin{Pin{"inArray", BOOLEAN}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			arr, ok := in[1].([]interface{})
			if !ok {
				out[0] = NewError("inArray requries array")
				return nil
			}
			for _, x := range arr {
				if x == in[0] {
					out[0] = true
					return nil
				}
			}
			out[0] = false
			return nil
		},
	}
}

// HasField returns true if the supplied string is a field of the supplied object
func HasField() Spec {
	return Spec{
		Name:    "hasField",
		Inputs:  []Pin{Pin{"field", STRING}, Pin{"object", OBJECT}},
		Outputs: []Pin{Pin{"hasField", BOOLEAN}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			obj, ok := in[1].(map[string]interface{})
			if !ok {
				out[0] = NewError("HasField requries map for object")
				return nil
			}
			field, ok := in[0].(string)
			if !ok {
				out[0] = NewError("HasField requires string for field")
				return nil
			}
			_, out[0] = obj[field]
			return nil
		},
	}
}

// InString returns true if substring is contained within string
func InString() Spec {
	return Spec{
		Name:    "inString",
		Inputs:  []Pin{Pin{"substring", STRING}, Pin{"string", STRING}},
		Outputs: []Pin{Pin{"inString", BOOLEAN}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			substring, ok := in[0].(string)
			if !ok {
				out[0] = NewError("inString requires string for substring")
				return nil
			}
			superstring, ok := in[1].(string)
			if !ok {
				out[0] = NewError("inString requires string for string")
				return nil
			}
			out[0] = strings.Contains(superstring, substring)
			return nil
		},
	}
}

// HasPrefix returns true if substring is prefix of string
func HasPrefix() Spec {
	return Spec{
		Name:    "hasPrefix",
		Inputs:  []Pin{Pin{"substring", STRING}, Pin{"string", STRING}},
		Outputs: []Pin{Pin{"inString", BOOLEAN}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			substring, ok := in[0].(string)
			if !ok {
				out[0] = NewError("HasPrefix requires strings")
				return nil
			}
			superstring, ok := in[1].(string)
			if !ok {
				out[0] = NewError("HasPrefix requires strings")
				return nil
			}
			out[0] = strings.HasPrefix(superstring, substring)
			return nil
		},
	}
}

// HasSuffix returns true if substring is prefix of string
func HasSuffix() Spec {
	return Spec{
		Name:    "hasSuffix",
		Inputs:  []Pin{Pin{"substring", STRING}, Pin{"string", STRING}},
		Outputs: []Pin{Pin{"inString", BOOLEAN}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			substring, ok := in[0].(string)
			if !ok {
				out[0] = NewError("HasSuffix requires strings")
				return nil
			}
			superstring, ok := in[1].(string)
			if !ok {
				out[0] = NewError("HasSuffix requires strings")
				return nil
			}
			out[0] = strings.HasSuffix(superstring, substring)
			return nil
		},
	}
}

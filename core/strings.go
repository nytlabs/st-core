package core

import "strings"

func StringConcat() Spec {
	return Spec{
		Name:    "concat",
		Inputs:  []Pin{Pin{"a", STRING}, Pin{"b", STRING}},
		Outputs: []Pin{Pin{"a+b", STRING}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			a, ok := in[0].(string)
			b, ok := in[1].(string)
			if !ok {
				out[0] = NewError("concat requires string")
				return nil
			}
			out[0] = a + b
			return nil
		},
	}
}

func StringSplit() Spec {
	return Spec{
		Name:    "split",
		Inputs:  []Pin{Pin{"a<sep>b", STRING}, Pin{"sep", STRING}},
		Outputs: []Pin{Pin{"[a,b]", ARRAY}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			ab, ok := in[0].(string)
			sep, ok := in[1].(string)
			if !ok {
				out[0] = NewError("split requires string")
				return nil
			}
			seperatedStrings := strings.Split(ab, sep)
			sepI := make([]interface{}, len(seperatedStrings))
			for i, v := range seperatedStrings {
				sepI[i] = v
			}
			out[0] = sepI
			return nil
		},
	}
}

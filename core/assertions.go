package core

import (
	"encoding/json"
	"strconv"
)

func IsError() Spec {
	return Spec{
		Name:    "isError",
		Inputs:  []Pin{Pin{"in", ANY}},
		Outputs: []Pin{Pin{"out", BOOLEAN}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			_, ok := in[0].(*stcoreError)
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
		Inputs:  []Pin{Pin{"in", ANY}},
		Outputs: []Pin{Pin{"out", BOOLEAN}},
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
		Inputs:  []Pin{Pin{"in", ANY}},
		Outputs: []Pin{Pin{"out", BOOLEAN}},
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
		Inputs:  []Pin{Pin{"in", ANY}},
		Outputs: []Pin{Pin{"out", BOOLEAN}},
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
		Inputs:  []Pin{Pin{"in", ANY}},
		Outputs: []Pin{Pin{"out", BOOLEAN}},
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
		Inputs:  []Pin{Pin{"in", ANY}},
		Outputs: []Pin{Pin{"out", BOOLEAN}},
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

func ToString() Spec {
	return Spec{
		Name:    "toString",
		Inputs:  []Pin{Pin{"in", ANY}},
		Outputs: []Pin{Pin{"out", STRING}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			switch t := in[0].(type) {
			case float64:
				out[0] = strconv.FormatFloat(t, 'f', -1, 64)
			case bool:
				out[0] = strconv.FormatBool(t)
			default:
				j, e := json.Marshal(t)
				if e != nil {
					out[0] = NewError(e.Error())
				}
				out[0] = string(j)
			}
			return nil
		},
	}
}

func ToNumber() Spec {
	return Spec{
		Name:    "toNumber",
		Inputs:  []Pin{Pin{"in", ANY}},
		Outputs: []Pin{Pin{"out", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			switch t := in[0].(type) {
			case float64:
				out[0] = t
			case bool:
				if t == true {
					out[0] = 1.0
				} else {
					out[1] = 0.0
				}
			case string:
				f, err := strconv.ParseFloat(t, 64)
				if err != nil {
					out[0] = NewError(err.Error())
				}
				out[0] = f
			default:
				out[0] = NewError("could not convert msg to float")
			}
			return nil
		},
	}
}

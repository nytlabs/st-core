package core

import "encoding/json"

func ParseJSON() Spec {
	return Spec{
		Name:    "parseJSON",
		Inputs:  []Pin{Pin{"in", STRING}},
		Outputs: []Pin{Pin{"out", ANY}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			msgstring, ok := in[0].(string)
			if !ok {
				out[0] = NewError("ParseJSON needs string")
				return nil
			}
			msgbytes := []byte(msgstring)
			var msg interface{}
			err := json.Unmarshal(msgbytes, &msg)
			if err != nil {
				out[0] = err
				return nil
			}
			out[0] = msg
			return nil
		},
	}
}

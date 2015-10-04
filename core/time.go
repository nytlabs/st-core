package core

import "time"

func Timestamp() Spec {
	return Spec{
		Name:    "timestamp",
		Inputs:  []Pin{Pin{"trigger", ANY}},
		Outputs: []Pin{Pin{"timestamp", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			out[0] = float64(time.Now().UnixNano() / 1000000)
			return nil
		},
	}
}

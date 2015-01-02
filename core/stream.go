package core

import "sync"

type Stream struct {
	sync.Mutex
	Out chan interface{}
	In  chan interface{}
}

func FromStream() Spec {
	return Spec{
		Name: "fromStream",
		Outupts: []Pin{
			Pin{"out"},
		},
		Shared: STREAM,
		Kernel: func(in, out, internal MessageMap, s Store, i chan Interrupt) Interrupt {
			stream := s.(*Stream)
			select {
			case out[0] = <-s.Out:
			case f := <-i:
				return f
			}
			return nil
		},
	}
}

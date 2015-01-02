package core

import (
	"log"
	"sync"
)

// Stream represents a two-way communication channel with the world outside of streamtools
type Stream struct {
	sync.Mutex
	Out  chan interface{}
	In   chan interface{}
	quit chan bool
}

func (s *Stream) Stop() {
	s.quit <- true
}

func NewStream() Store {
	stream := &Stream{
		Out:  make(chan interface{}),
		In:   make(chan interface{}),
		quit: make(chan bool),
	}
	// closure that interfaces with whatever tech
	go func() {
		outmessage := map[string]string{"hello": "world"}
		for {
			select {
			case stream.Out <- outmessage:
			case <-stream.quit:
				return
			case m := <-stream.In:
				log.Println(m)
			}
		}
	}()
	return stream
}

func streamReceive() Spec {
	return Spec{
		Name: "streamReceive",
		Outputs: []Pin{
			Pin{"out"},
		},
		Shared: STREAM,
		Kernel: func(in, out, internal MessageMap, s Store, i chan Interrupt) Interrupt {
			stream := s.(*Stream)
			select {
			case out[0] = <-stream.Out:
			case f := <-i:
				return f
			}
			return nil
		},
	}
}

func streamSend() Spec {
	return Spec{
		Name: "streamSend",
		Inputs: []Pin{
			Pin{"in"},
		},
		Shared: STREAM,
		Kernel: func(in, out, internal MessageMap, s Store, i chan Interrupt) Interrupt {
			stream := s.(*Stream)
			select {
			case stream.In <- in[0]:
			case f := <-i:
				return f
			}
			return nil
		},
	}
}

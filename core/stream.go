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

// Stop causes the stream to close. No more messages will show up, and no more messages can be sent
func (s *Stream) Stop() {
	s.quit <- true
}

// NewStream returns a new stream Store.
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

// StreamRecieve receives messages from the Stream store.
//
// OutPin 0: received message
func StreamReceive() Spec {
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

// StreamSend publishes the inbound message to the Stream store.
//
// InPin 0: message to send
func StreamSend() Spec {
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

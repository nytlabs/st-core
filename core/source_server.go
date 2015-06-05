package core

import "sync"

func Server() SourceSpec {
	return SourceSpec{
		Name: "Server",
		Type: SERVER,
		New:  NewServer,
	}
}

type Server struct {
	quit chan bool
	Out  chan Message // this channel is used by any block that would like to receive messages
	sync.Mutex
}

func (s Server) GetType() SourceType {
	return SERVER
}

func (s *NSQ) SetSourceParameter(name, value string) {
	switch name {
	}
}

func (s *NSQ) Describe() map[string]string {
	return map[string]string{}
}

func NewServer() Source {
	out := make(chan Message)
	server := &Server{
		quit: make(chan bool),
		Out:  out,
	}
	return server
}

func (s Server) Serve() {
}

func (s Server) Stop() {
	s.quit <- true
}

// NSQRecieve receives messages from the NSQ system.
//
// OutPin 0: received message
func FromRequest() Spec {
	return Spec{
		Name: "FromRequest",
		Outputs: []Pin{
			Pin{"out", OBJECT},
		},
		Source: SERVER,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			server := s.(*Server)
			select {
			case out[0] = <-stream.Out:
			case f := <-i:
				return f
			}
			return nil
		},
	}
}

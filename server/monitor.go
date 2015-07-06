package server

import "github.com/nytlabs/st-core/core"

// code for inferring the current state of the st-core pattern
// and emitting it over the websocket

func (s *Server) MonitorMux(id int, c chan core.MonitorMessage) {
	for m := range c {
		s.websocketBroadcast(Update{Action: INFO, Type: BLOCK, Data: wsInfo{wsId{id}, m}})
	}
}

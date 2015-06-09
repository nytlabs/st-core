package server

import (
	"time"

	"github.com/nytlabs/st-core/core"
)

// code for inferring the current state of the st-core pattern
// and emitting it over the websocket

func (s *Server) MonitorMux(id int, c chan time.Time) {
	for m := range c {
		s.websocketBroadcast(Update{Action: UPDATE, Type: BLOCK, Data: wsBlock{wsAlert{wsId{id}, core.CRANK, m}}})
	}
}

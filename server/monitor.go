package server

import (
	"time"

	"github.com/nytlabs/st-core/core"
)

// code for inferring the current state of the st-core pattern
// and emitting it over the websocket

func (s *Server) MonitorMux(id int, c chan core.MonitorMessage, query chan struct{}, quit chan struct{}) {
	expire := time.NewTimer(time.Duration(250 * time.Millisecond))
	var state core.MonitorMessage
	running := false
	for {
		select {
		case m := <-c:
			state = m
			expire.Reset(time.Duration(250 * time.Millisecond))
			if !running {
				running = true
				s.websocketBroadcast(Update{Action: INFO, Type: BLOCK, Data: wsInfo{wsId{id}, core.MonitorMessage{
					core.BI_RUNNING,
					nil,
					//					time.Now(),
				}}})
			}
		case <-expire.C:
			s.websocketBroadcast(Update{Action: INFO, Type: BLOCK, Data: wsInfo{wsId{id}, state}})
			running = false
		case <-quit:
			return
		case <-query:
			if running {
				s.websocketBroadcast(Update{Action: INFO, Type: BLOCK, Data: wsInfo{wsId{id}, core.MonitorMessage{
					core.BI_RUNNING,
					nil,
					//					time.Now(),
				}}})
			} else {
				s.websocketBroadcast(Update{Action: INFO, Type: BLOCK, Data: wsInfo{wsId{id}, state}})
			}
		}
	}
}

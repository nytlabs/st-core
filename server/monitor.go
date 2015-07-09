package server

import (
	"time"

	"github.com/nytlabs/st-core/core"
)

// monitor infers the current state of an st-core block and emits it over
// websocket.
//
// every start of a block in server also starts a monitor routine. each pin in
// a block, whether in a receive or broadcast state, emits a
// core.MonitorMessage to the monitor. A kernel process also emits a message to
// the monitor. The monitor's job is to only emit state change, i.e., going
// from one blocking state to another.
//
// currently, the sends from block are blocking, meaning that they are limtied
// by the speed in which monitor can process messages.
//
// TODO: test the overhead of blocking sends from block to monitor. currently,
// every block ges a monitor by default. an optimizing step may be to disable
// it by default, and affording the querying of state by some API handle.
//
// in that capacity, the state should live in core.Block, and MonitorMux would
// simply use a mutex to transmit it, instead of any continous channel traffic
// from block to Monitor

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
				}}})
			} else {
				s.websocketBroadcast(Update{Action: INFO, Type: BLOCK, Data: wsInfo{wsId{id}, state}})
			}
		}
	}
}

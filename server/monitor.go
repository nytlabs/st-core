package server

import (
	"log"

	"github.com/nytlabs/st-core/core"
)

// code for inferring the current state of the st-core pattern
// and emitting it over the websocket

func (s *Server) MonitorMux(id int, c chan core.BlockAlert) {

	key := map[core.BlockAlert]string{
		core.BLOCKED:   "blocked",
		core.UNBLOCKED: "unblocked",
	}

	for {
		m := <-c
		log.Println(id, key[m])
		s.websocketBroadcast(Update{Action: UPDATE, Type: ALERT, Data: wsAlert{
			Id:    id,
			Alert: key[m],
		}})
	}

}

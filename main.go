package main

import (
	"time"

	"github.com/nytlabs/st-core/core"
	"github.com/thejerf/suture"
)

func main() {
	supervisor := suture.NewSimple("st-core")

	log := core.NewBlock(core.Library["log"])
	delay := core.NewBlock(core.Library["delay"])
	pusher := core.NewBlock(core.Library["pusher"])

	supervisor.Add(log)
	supervisor.Add(delay)
	supervisor.Add(pusher)

	pusher.GetOutput("out").Connect(delay.GetInput("in"))
	delay.GetOutput("out").Connect(log.GetInput("in"))

	timer1 := time.NewTimer(5 * time.Second)

	supervisor.ServeBackground()

	<-timer1.C
}

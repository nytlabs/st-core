package main

import (
	"time"

	"github.com/nytlabs/st-core/blocks"
	"github.com/thejerf/suture"
)

func main() {
	supervisor := suture.NewSimple("st-core")

	p := blocks.NewPusher()
	d := blocks.NewDelay()
	l := blocks.NewLog("logger")

	supervisor.Add(p)
	supervisor.Add(d)
	supervisor.Add(l)

	p.Connect("out", d.GetInput("in"))
	d.Connect("out", l.GetInput("in"))

	timer := time.NewTimer(500 * time.Millisecond)

	go supervisor.ServeBackground()

	<-timer.C
	supervisor.Stop()
}

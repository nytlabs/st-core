package main

import (
	"time"

	"github.com/nytlabs/st-core/core"
	"github.com/thejerf/suture"
)

func main() {
	supervisor := suture.NewSimple("st-core")

	p := core.NewPusher()
	d := core.NewDelay()
	l := core.NewLog("logger")

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

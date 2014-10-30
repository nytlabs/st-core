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

	ptoken := supervisor.Add(p)
	_ = supervisor.Add(d)
	_ = supervisor.Add(l)

	p.Connect("out", d.GetInput("in"))
	d.Connect("out", l.GetInput("in"))

	timer1 := time.NewTimer(500 * time.Millisecond)
	timer2 := time.NewTimer(2 * time.Second)

	go supervisor.ServeBackground()

	<-timer1.C
	supervisor.Remove(ptoken)

	<-timer2.C
	supervisor.Stop()
}

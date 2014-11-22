package main

import (
	"time"

	"github.com/nytlabs/st-core/core"
	"github.com/thejerf/suture"
)

func main() {
	supervisor := suture.NewSimple("st-core")

	p1 := core.NewPusher("testPusher_a")
	p2 := core.NewPusher("testPusher_b")

	a := core.NewPlus("plus")
	d := core.NewDelay("delay")
	l := core.NewLog("log")

	supervisor.Add(a)
	supervisor.Add(d)
	supervisor.Add(l)
	supervisor.Add(p1)
	supervisor.Add(p2)

	p1.Connect("out", a.GetInput("addend 1"))
	p2.Connect("out", a.GetInput("addend 2"))
	a.Connect("out", d.GetInput("in"))
	d.Connect("out", l.GetInput("in"))

	timer1 := time.NewTimer(5 * time.Second)

	supervisor.ServeBackground()

	<-timer1.C
}

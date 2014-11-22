package main

import (
	"log"
	"time"

	"github.com/nytlabs/st-core/core"
	"github.com/thejerf/suture"
)

func main() {
	supervisor := suture.NewSimple("st-core")
	supervisor.ServeBackground()

	p := core.NewPusher("testPusher")
	err := p.Inputs["value"].SetValue(2.1)
	if err != nil {
		log.Fatal(err)
	}
	ptoken := supervisor.Add(p)

	f := core.NewF("f")
	g := core.NewG("g")

	supervisor.Add(f)
	supervisor.Add(g)

	p.Connect("out", f.GetInput("in"))
	f.Connect("out", g.GetInput("in"))

	d := core.NewDelay("delay")
	l := core.NewLog("logger")

	supervisor.Add(d)
	supervisor.Add(l)

	g.Connect("out", d.GetInput("in"))
	d.Connect("out", l.GetInput("in"))

	timer1 := time.NewTimer(3 * time.Second)
	timer2 := time.NewTimer(9 * time.Second)

	<-timer1.C
	log.Println("merging f and g")

	h := f.Merge(g)

	<-timer2.C
	supervisor.Stop()
}

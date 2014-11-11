package main

import (
	"time"

	"github.com/nytlabs/st-core/core"
	"github.com/nytlabs/st-core/stores"
	"github.com/thejerf/suture"
)

func main() {
	supervisor := suture.NewSimple("st-core")

	blocks := make(map[string]suture.ServiceToken)

	p := core.NewPusher("testPusher")
	d := core.NewDelay("delay")
	l := core.NewLog("logger")
	s := stores.NewKeyValue("store")
	sSet := stores.NewKeyValueSet("setter")
	sGet := stores.NewKeyValueGet("getter")

	blocks[p.Name] = supervisor.Add(p)

	_ = supervisor.Add(d)
	_ = supervisor.Add(l)
	_ = supervisor.Add(s)
	_ = supervisor.Add(sSet)
	_ = supervisor.Add(sGet)

	p.Connect("out", d.GetInput("in"))
	d.Connect("out", l.GetInput("in"))

	sSet.ConnectStore(s)
	sGet.ConnectStore(s)

	timer1 := time.NewTimer(2 * time.Second)
	timer2 := time.NewTimer(5 * time.Second)

	supervisor.ServeBackground()

	<-timer1.C
	supervisor.Remove(blocks["testPusher"])

	<-timer2.C
	supervisor.Stop()
}

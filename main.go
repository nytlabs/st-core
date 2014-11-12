package main

import (
	"log"
	"time"

	"github.com/nytlabs/st-core/core"
	"github.com/thejerf/suture"
)

func main() {
	supervisor := suture.NewSimple("st-core")

	blocks := make(map[string]suture.ServiceToken)

	/*
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
	*/

	m := core.NewMap("test1")

	testValue := map[string]interface{}{
		"hello": "world",
	}

	err := m.Inputs["in"].SetValue(testValue)
	if err != nil {
		log.Fatal(err)
	}

	testMapping := map[string]interface{}{
		"foo": ".hello",
	}

	err = m.Inputs["mapping"].SetValue(testMapping)
	if err != nil {
		log.Fatal(err)
	}

	l := core.NewLog("logger")
	d := core.NewDelay("delay")

	supervisor.Add(m)
	supervisor.Add(l)
	supervisor.Add(d)

	m.Connect("out", d.GetInput("in"))
	d.Connect("out", l.GetInput("in"))

	timer1 := time.NewTimer(2 * time.Second)
	timer2 := time.NewTimer(5 * time.Second)

	supervisor.ServeBackground()

	<-timer1.C
	supervisor.Remove(blocks["testPusher"])

	<-timer2.C
	supervisor.Stop()
}

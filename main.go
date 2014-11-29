package main

import (
	"time"

	"github.com/nytlabs/st-core/core"
	"github.com/thejerf/suture"
)

func main() {
	supervisor := suture.NewSimple("st-core")
	supervisor.ServeBackground()

	log := core.NewBlock(core.Library["log"])
	delay := core.NewBlock(core.Library["delay"])
	pusher := core.NewBlock(core.Library["pusher"])
	set := core.NewBlock(core.Library["set"])

	pusher.GetInput("value").SetValue("bob")
	set.GetInput("key").SetValue("name")
	delay.GetInput("duration").SetValue("1s")

	supervisor.Add(log)
	supervisor.Add(delay)
	_ = supervisor.Add(pusher)
	setToken := supervisor.Add(set)

	pusher.GetOutput("out").Connect(delay.GetInput("in"))
	delay.GetOutput("out").Connect(set.GetInput("in"))
	set.GetOutput("out").Connect(log.GetInput("in"))

	timer1 := time.NewTimer(6 * time.Second)
	timer2 := time.NewTimer(9 * time.Second)

	<-timer1.C

	supervisor.Remove(setToken)

	<-timer2.C
	supervisor.Stop()
}

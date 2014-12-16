package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nikhan/go-fetch"
	"github.com/nytlabs/st-core/core"
	"github.com/thejerf/suture"
)

func main() {
	supervisor := suture.NewSimple("st-core")
	supervisor.ServeBackground()

	library := core.GetLibrary()

	b := core.NewBlock(library["plus"])
	d := core.NewBlock(library["delay"])
	l := core.NewBlock(library["log"])
	o := core.NewBlock(library["set"])
	l2 := core.NewBlock(library["log"])

	a := supervisor.Add(b)
	_ = supervisor.Add(d)
	_ = supervisor.Add(l)
	_ = supervisor.Add(o)
	_ = supervisor.Add(l2)

	b.Connect(0, d.Input(0).C)
	d.Connect(0, o.Input(1).C)
	o.Connect(0, l.Input(0).C)

	path, _ := fetch.Parse(".test")
	fmt.Println(path)

	b.RouteValue(core.RouteID(0), 1.1)
	d.RouteValue(core.RouteID(1), "10ms")
	o.RouteValue(core.RouteID(0), "test")

	go func() {
		for {
			b.RouteValue(core.RouteID(1), rand.Float64()*10.0)
			time.Sleep(time.Duration(rand.Intn(2)+1) * time.Millisecond)
		}

	}()

	time.Sleep(500 * time.Millisecond)

	fmt.Println("Disconnected!")

	d.Disconnect(0, o.Input(1).C)

	l2.RoutePath(core.RouteID(0), path)
	o.Connect(core.RouteID(0), l2.Input(0).C)

	time.Sleep(300 * time.Millisecond)

	d.Connect(0, o.Input(1).C)
	fmt.Println("Connected")

	time.Sleep(100 * time.Millisecond)

	supervisor.Remove(a)

	supervisor.Stop()
}

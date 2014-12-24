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

	set := core.NewBlock(library["set"])
	log := core.NewBlock(library["log"])
	kvset := core.NewBlock(library["kvSet"])
	kvdump := core.NewBlock(library["kvDump"])
	delay := core.NewBlock(library["delay"])
	setdelay := core.NewBlock(library["delay"])

	kv := core.NewKeyValue()

	_ = supervisor.Add(set)
	_ = supervisor.Add(log)
	_ = supervisor.Add(kvset)
	_ = supervisor.Add(kvdump)
	_ = supervisor.Add(delay)
	_ = supervisor.Add(setdelay)

	path, _ := fetch.Parse(".test")
	set.RouteValue(0, "test")
	kvset.RouteValue(1, "whatever")
	kvset.RoutePath(0, path)
	kvset.Store(kv)
	set.Connect(0, kvset.Input(0).C)
	delay.RouteValue(0, "bang")
	delay.RouteValue(1, "1s")
	delay.Connect(0, kvdump.Input(0).C)
	kvdump.Connect(0, log.Input(0).C)
	kvset.Connect(0, setdelay.Input(0).C)
	setdelay.RouteValue(1, "100ms")
	setdelay.Connect(0, log.Input(0).C)
	kvdump.Store(kv)
	go func() {
		for {
			set.RouteValue(1, fmt.Sprintf("%2.2f", rand.Float64()*10.0))
			time.Sleep(time.Duration(rand.Intn(2)+1) * time.Millisecond)
		}
	}()

	time.Sleep(10 * time.Second)

	supervisor.Stop()

	/*
		s := core.NewServer()
		r := mux.NewRouter()
		r.HandleFunc("/", s.WebsocketHandler).Methods("GET")
		r.HandleFunc("/group", s.GetGroupHandler).Methods("GET")
		r.HandleFunc("/group", s.CreateGroupHandler).Methods("POST")
		r.HandleFunc("/block", s.CreateBlockHandler).Methods("POST")
		r.HandleFunc("/connections", s.CreateConnectionHandler).Methods("POST")
		http.Handle("/", r)

		log.Println("serving on 7071")
		err := http.ListenAndServe(":7071", nil)
		if err != nil {
			log.Panicf(err.Error())
		}
	*/

	/*
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
	*/
}

package main

import (
	"fmt"
	"time"

	"github.com/nytlabs/st-core/core"
)

func main() {
	sink := make(chan core.Message)
	set := core.NewBlock(core.GetLibrary()["set"])
	go set.Serve()

	set.SetRoute(0, "test")
	set.SetRoute(1, "ok")

	set.Connect(0, sink)

	last := set
	for i := 0; i < 1000; i++ {
		last.Disconnect(0, sink)
		newI := core.NewBlock(core.GetLibrary()["identity"])
		go newI.Serve()
		nn, _ := newI.GetRoute(0)
		last.Connect(0, nn.C)
		newI.Connect(0, sink)
		last = newI

		start := time.Now().UnixNano()
		for i := 0; i < 1000; i++ {
			_ = <-sink
		}
		end := time.Now().UnixNano()
		x := time.Duration(time.Duration(end-start) * time.Nanosecond).Seconds()
		fmt.Println(fmt.Sprintf("%d identity blocks", i+1), 1000.0/x, "msgs/sec")
	}
}

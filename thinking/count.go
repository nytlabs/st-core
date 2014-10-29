package main

import (
	"fmt"
	"log"

	"github.com/thejerf/suture"
)

type Count struct {
	count    int
	inChan   chan interface{}
	outChan  chan interface{}
	pollChan chan interface{}
	stopChan chan interface{}
}

func (b *Count) Stop() {
	fmt.Println("stopping")
	b.stopChan <- true
}

func (b *Count) Serve() {
	for {
		select {
		case <-b.inChan:
			b.count++
		case <-b.pollChan:
			b.outChan <- b.count
		case <-b.stopChan:
			log.Println("stopping")
			return
		}
	}
}

func main() {
	supervisor := suture.NewSimple("st-core")
	b := &Count{
		inChan:   make(chan interface{}),
		outChan:  make(chan interface{}),
		pollChan: make(chan interface{}),
	}
	supervisor.Add(b)
	go supervisor.ServeBackground()
	b.inChan <- true
	b.inChan <- true
	go func(b *Count) { b.pollChan <- true }(b)
	log.Println(<-b.outChan)
	log.Println(supervisor.Services())
	supervisor.Stop()
}

package main

import (
	"fmt"
	"log"
)

type Plus struct {
	leftTermChan  chan interface{}
	rightTermChan chan interface{}
	resultChan    chan interface{}
	stopChan      chan interface{}
}

func (b *Plus) Stop() {
	fmt.Println("stopping")
	b.stopChan <- true
}

func (b *Plus) Serve() {

	left := make(chan interface{})
	right := make(chan interface{})

	go func() {
		for {
			x := <-left
			y := <-right
			xf := x.(float64)
			yf := y.(float64)
			b.resultChan <- xf + yf
		}
	}()

	for {
		select {
		case leftTerm := <-b.leftTermChan:
			go func() { left <- leftTerm }()
		case rightTerm := <-b.rightTermChan:
			go func() { right <- rightTerm }()
		case <-b.stopChan:
			log.Println("stopping")
			return
		}
	}
}

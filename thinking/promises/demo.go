package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Message struct {
	id int
}

type Promise chan Message

type Signal chan Promise

type Block struct {
	concurrency int
	kernel      func(Message) Message
}

func NewBlock(concurrency int, kernel func(Message) Message) *Block {
	return &Block{concurrency: concurrency, kernel: kernel}
}

func (b *Block) Pipe(in Signal) Signal {
	out := make(Signal, b.concurrency)

	go func() {
		defer close(out)
		for inPromise := range in {

			outPromise := make(Promise)
			out <- outPromise

			go func() {
				m := <-inPromise
				outPromise <- b.kernel(m)
			}()
		}
	}()
	return out
}

func Boundary(in Signal) chan Message {
	messages := make(chan Message)
	go func() {
		defer close(messages)
		for promise := range in {
			messages <- <-promise
		}
	}()
	return messages
}

func main() {

	latentBlock1 := NewBlock(10, func(m Message) Message {
		sleep := rand.Intn(10) + 1
		time.Sleep(time.Millisecond * time.Duration(sleep))
		fmt.Printf("kernel1: %d\n", m.id)
		return m
	})

	startSignal := make(Signal)
	endSignal := latentBlock1.Pipe(startSignal)

	go func() {
		for i := 0; i < 100; i++ {
			startPromise := make(Promise)
			startSignal <- startPromise
			startPromise <- Message{id: i}
		}
		close(startSignal)
	}()

	synchronizedMessages := Boundary(endSignal)

	for m := range synchronizedMessages {
		fmt.Printf("sync: %d\n", m.id)
	}
}

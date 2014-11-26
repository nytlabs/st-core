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

	// concurrency controls the total number of promises that can be
	// outstanding at any one time.
	concurrency int

	// kernel simplifies message processing to be one-to-one mesage in message
	// out.
	kernel func(Message) Message
}

func NewBlock(concurrency int, kernel func(Message) Message) *Block {
	return &Block{concurrency: concurrency, kernel: kernel}
}

// Pipe directs a new stream of input to this block and returns the output
// Signal.
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

// Bundary converts a Signal to a channel of ordered Messages.
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

	latencyBlock := NewBlock(10, func(m Message) Message {

		// Sleep to simulate a latency-heavy task, such as a web request or DB
		// request.
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)+1))

		// Anounce that the task is complete for a given message and pass it
		// on to the next block.
		fmt.Printf("kernel: %d\n", m.id)
		return m
	})

	// Create a new signal and pipe it to the latency-simulating block.
	startSignal := make(Signal)
	endSignal := latencyBlock.Pipe(startSignal)

	// Push some in-order dummy messages into the system. The only purpose of
	// the ID is to later identify if the sequence is still in order.
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

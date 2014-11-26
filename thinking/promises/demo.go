package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Signal chan chan Message

type Message struct {
	id int
}

type Block struct {
	c      int
	kernel func(Message) Message
}

func NewBlock(c int, kernel func(Message) Message) *Block {
	return &Block{c: c, kernel: kernel}
}

func (b *Block) Pipe(in Signal) Signal {
	out := make(Signal)
	go func() {
		for {
			inPromise, ok := <-in
			if !ok {
				close(out)
				return
			}

			outPromise := make(chan Message)
			out <- outPromise

			go func() {
				m := <-inPromise
				outPromise <- b.kernel(m)
			}()
		}
	}()
	return out
}

func main() {

	latentBlock := NewBlock(500, func(m Message) Message {
		sleep := rand.Intn(10) + 1
		time.Sleep(time.Millisecond * time.Duration(sleep))
		fmt.Printf("kernel: %d\n", m.id)
		return m
	})

	startSignal := make(Signal)
	endSignal := latentBlock.Pipe(startSignal)

	go func() {
		for i := 0; i < 100; i++ {
			startPromise := make(chan Message)
			startSignal <- startPromise
			startPromise <- Message{id: i}
		}
		close(startSignal)
	}()

	fmt.Println(latentBlock)

	endPromises := make([]chan Message, 0)
	for {
		endPromise, ok := <-endSignal
		if !ok {
			break
		}
		endPromises = append(endPromises, endPromise)
	}

	for i, endPromise := range endPromises {
		m := <-endPromise
		fmt.Printf("sync: %d, expecting %d\n", m.id, i)
	}
}

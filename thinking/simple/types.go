package main

import (
	"github.com/nikhan/go-fetch"
	"sync"
)

type RouteID int
type Message interface{}
type Connection chan Message
type MessageMap map[RouteID]Message
type InterruptFunc func() bool
type KernelFunc func(MessageMap, MessageMap, chan InterruptFunc) InterruptFunc

type Pin struct {
	Name string
}

type Spec struct {
	Inputs  []Pin
	Outputs []Pin
	Kernel  KernelFunc
}

type Route struct {
	Name  string
	Path  *fetch.Query
	Value *Message
	C     chan Message
}

type Output struct {
	Name        string
	Connections map[Connection]struct{}
}

type BlockState struct {
	inputValues  MessageMap
	outputValues MessageMap
}

type BlockRouting struct {
	Inputs        []Route
	Outputs       []Output
	InterruptChan chan InterruptFunc
	sync.RWMutex
}

type Block struct {
	state   BlockState
	routing BlockRouting
	kernel  KernelFunc
}

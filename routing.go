//
// this is a test to try to create a thread-safe set of outputs
//
// a router is a collection of inchans from downstream blocks
// a single output is a router
// a block has a collection of outputs
//
// this lets you add and remove outputs dynamically in a thread-safe
// manner as well as connections to each of those outputs
//
// when you call broadcast on a single output, it sends a message to every connection
// attached to that output
//

package main

import(
    "fmt"
    "sync"
)

type Route chan interface{}

func NewRouter() *Router{
    return &Router{
        Routes: make(map[Route]bool),
    }
}

type Router struct {
    sync.RWMutex
    Routes map[Route]bool
}

func (rh *Router) Add(r Route) bool {
    rh.Lock()
    defer rh.Unlock()
    _, ok := rh.Routes[r]
    if ok {
        return false
    }
    rh.Routes[r] = true
    return true
}

func (rh *Router) Remove(r Route) bool {
    rh.Lock()
    defer rh.Unlock()
    _, ok := rh.Routes[r]
    if !ok {
        return false
    }
    delete(rh.Routes, r)
    return true
}

func (rh *Router) Broadcast(m interface{}) {
    rh.RLock()
    for r, _ := range rh.Routes {
        r <- m
    }
    rh.RUnlock()
}

func NewBlock() *Block{
    return &Block{
        Outputs:make(map[string]*Router),
    }
}

type Block struct {
    Outputs map[string]*Router
    sync.RWMutex
}

func (b *Block) Add(id string) bool {
    b.Lock()
    defer b.Unlock()
    _, ok := b.Outputs[id]
    if ok {
        return false
    }
    b.Outputs[id] = NewRouter()
    return true
}

func (b *Block) Remove(id string) bool {
    b.Lock()
    defer b.Unlock()
    _, ok := b.Outputs[id]
    if !ok {
        return false
    }
    delete(b.Outputs, id)
    return true
}

func(b *Block) Connect(id string, r Route) bool {
    b.RLock()
    defer b.RUnlock()
    ok := b.Outputs[id].Add(r)
    return ok
}

func(b *Block) Disconnect(id string, r Route) bool {
    b.RLock()
    defer b.RUnlock()
    ok := b.Outputs[id].Remove(r)
    return ok
}

func(b *Block) Broadcast(id string, m interface{}) {
    b.RLock()
    defer b.RUnlock()
    b.Outputs[id].Broadcast(m)
}

func main(){
    connections := [10]Route{}

    for i, _ := range connections {
        connections[i] = make(Route, 1)
    }

    r := NewBlock()

    r.Add("test") // add output "test" to the block

    for _, c := range connections {
        r.Connect("test", c) // add 10 dummy connections to that output
    }

    r.Broadcast("test", "hello") // broadcast "hello" to all connections on output "test"

    // 10 hellos
    for _, c := range connections {
        m := <- c
        fmt.Println(m)
    }
}
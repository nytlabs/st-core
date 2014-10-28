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
    "time"
)

type Route chan interface{}

func NewRouter() *Router{
    return &Router{
        Routes: make(map[Route]bool),
    }
}

type Router struct {
    sync.Mutex
    Routes map[Route]bool
}

func (rh *Router) Add(r Route) bool {
    //rh.Lock()
    //defer rh.Unlock()
    _, ok := rh.Routes[r]
    if ok {
        return false
    }
    rh.Routes[r] = true
    return true
}

func (rh *Router) Remove(r Route) bool {
    //rh.Lock()
    //defer rh.Unlock()
    _, ok := rh.Routes[r]
    if !ok {
        return false
    }
    delete(rh.Routes, r)
    return true
}

func (rh *Router) Broadcast(m interface{}) {
    rh.Lock()
    routes := rh.Routes
    rh.Unlock()
    for r, _ := range routes {
        r <- m
    }
}

func NewBlock() *Block{
    return &Block{
        Inputs:make(map[string]Route),
        Outputs:make(map[string]*Router),
    }
}

type Block struct {
    Inputs map[string]Route
    Outputs map[string]*Router
    sync.Mutex
}

func (b *Block) AddInput(id string) bool {
    b.Lock()
    defer b.Unlock()
    _, ok := b.Inputs[id]
    if ok {
        return false
    }
    b.Inputs[id] = make(Route)
    return true
}

func (b *Block) RemoveInput(id string) bool {
    b.Lock()
    defer b.Unlock()
    _, ok := b.Inputs[id]
    if !ok {
        return false
    }
    delete(b.Inputs, id)
    return true
}

func (b *Block) Input(id string) Route {
    b.Lock()
    defer b.Unlock()
    input, ok := b.Inputs[id]
    if !ok {
        return nil
    }
    return input
}

func (b *Block) Recieve(id string) Route {
    b.Lock()
    m := b.Inputs[id]
    b.Unlock()
    return m
}

func (b *Block) AddOutput(id string) bool {
    b.Lock()
    defer b.Unlock()
    _, ok := b.Outputs[id]
    if ok {
        return false
    }
    b.Outputs[id] = NewRouter()
    return true
}

func (b *Block) RemoveOutput(id string) bool {
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
    b.Lock()
    ok := b.Outputs[id].Add(r)
    b.Unlock()
    return ok
}

func(b *Block) Disconnect(id string, r Route) bool {
    b.Lock()
    ok := b.Outputs[id].Remove(r)
    b.Unlock()
    return ok
}

func(b *Block) Broadcast(id string, m interface{}) {
    b.Lock()
    r := b.Outputs[id]
    b.Unlock()
    r.Broadcast(m)
}


func Pusher() *Block{
    b := NewBlock()
    b.AddOutput("out")
    go func(){
        i := 0
        for{
            i++
            b.Broadcast("out",i)
        }
    }()
    return b
}

func Delay() *Block{
    b := NewBlock()
    b.AddInput("in")
    b.AddOutput("out")
    go func(){
        for{
            m := <- b.Recieve("in")
            time.Sleep(10 * time.Millisecond)
            b.Broadcast("out", m)
        }
    }()
    return b
}

func Log(name string) *Block{
    b := NewBlock()
    b.AddInput("in")
    go func(){
        for{
            m := <- b.Recieve("in")
            fmt.Println(name, ": ", m)
        }
    }()
    return b
}

func Plus() *Block{
    b := NewBlock()
    b.AddInput("addend 1")
    b.AddInput("addend 2")
    b.AddOutput("out")
    go func(){
        for{
            bdd := <- b.Recieve("addend 2")
            add := <- b.Recieve("addend 1")
            c := add.(int) + bdd.(int)
            fmt.Println(add, bdd)
            b.Broadcast("out", c)
        }
    }()

    return b
}


func main(){

    go func(){
        t := time.NewTicker(1 * time.Second)
        for{
            select{
            case <-t.C:
            }
        }
    }()

    p := Pusher()
    d := Delay()
    plus := Plus()
    l := Log("logger")

    time.Sleep(1 * time.Second)

    delayIn := d.Input("in")
    p.Connect("out", delayIn)
    plus_1 := plus.Input("addend 1")
    plus_2 := plus.Input("addend 2")
    d.Connect("out", plus_1)
    p.Connect("out", plus_2)
    login :=  l.Input("in")
    plus.Connect("out",login)

    <- make(chan bool)
    /*connections := [10]Route{}

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
    }*/
}
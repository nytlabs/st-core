package main

import (
    "fmt"
    "time"
)

type Delay struct{
    in chan interface{}
    out chan interface{}
    ping chan chan string
    quit chan chan bool
    workerDone chan bool
}

func (d *Delay) Serve(){
    go func(){
        for{
            for {
                m, ok := <- d.in
                if !ok {
                    d.workerDone <- true
                    return
                }
                time.Sleep(1 * time.Second)
                d.out <- m
            }
        }
    }()

    for{
        select {
        case m := <- d.ping:
            m <- "Delay"
        case q := <- d.quit:
            close(d.in)
            <-d.workerDone
            q <- true            
            return
        }
    }
}

type Pusher struct{
    in chan interface{}
    out chan interface{}
    ping chan chan string
    quit chan chan bool
    workerDone chan bool
}

func (d *Pusher) Serve(){
    quit := make(chan bool)
    go func(){
        m := 0
        for{
            select{
            case <-quit:
                d.workerDone <- true
                return
            default:
                m++ 
                d.out <- m
            }
        }
    }()

    for{
        select {
        case m := <- d.ping:
            m <- "Pusher"
        case q := <- d.quit:
            quit <- true
            <-d.workerDone
            q <- true            
            return
        }
    }
}

type Logger struct{
    in chan interface{}
    ping chan chan string
    quit chan chan bool
    workerDone chan bool
}

func (d *Logger) Serve(){
    go func(){
        for {
            m, ok := <- d.in
            if !ok {
                d.workerDone <- true
                return
            }
            fmt.Println(m)
        }
    }()

    for{
        select {
        case m := <- d.ping:
            m <- "Logger"
        case q := <- d.quit:
            close(d.in)
            <-d.workerDone
            q <- true            
            return
        }
    }
}

func main(){
    pusherToDelay := make(chan interface{})
    delayToLog := make(chan interface{})

    delayPing := make(chan chan string, 1)
    pusherPing := make(chan chan string, 1)
    loggerPing := make(chan chan string, 1)

    p := Pusher{
        in: make(chan interface{}),
        ping: pusherPing,
        out: pusherToDelay,
        quit: make(chan chan bool, 1),
        workerDone: make(chan bool, 1),
    }

    d := Delay{
        in: pusherToDelay,
        ping: delayPing,
        out: delayToLog,
        quit: make(chan chan bool, 1),
        workerDone: make(chan bool, 1),
    }

    l := Logger{
        in: delayToLog,
        ping: loggerPing,
        quit: make(chan chan bool, 1),
        workerDone: make(chan bool, 1),
    }

    go p.Serve()
    go d.Serve()
    go l.Serve()
    
    t:= time.NewTicker(200 * time.Millisecond)
    stop:= time.NewTimer(5 * time.Second)
    for{
        select {
        case <- t.C:
            respDelay := make(chan string)
            respPusher := make(chan string)
            respLogger := make(chan string)

            delayPing <- respDelay
            pusherPing <- respPusher
            loggerPing <- respLogger

            m1 := <- respDelay
            m2 := <- respPusher
            m3 := <- respLogger

            fmt.Println(m1, m2, m3)
        case <-stop.C:
            quitPusher := make(chan bool)
            //quitDelay := make(chan bool)
            //quitLogger := make(chan bool)
            p.quit <- quitPusher
            //d.quit <- quitDelay
            //l.quit <- quitLogger

            <- quitPusher
            fmt.Println("quit pusher")
            return
        }
    }
}

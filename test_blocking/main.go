package main

import (
    "fmt"
    "time"
)

type Pusher struct{
    Ping chan chan string
    out chan interface{}
}

func (p *Pusher) Run(){
    for {
        select{
            case response := <- p.Ping:
                response <- "pusher"
            default:
        }
        if len(p.out) == 0 {
            p.out <- ">>>>>>>>bang"
        }
    }
}

type Delay struct {
    Ping chan chan string
    in chan interface{} 
    out chan interface{}
}

func (d *Delay) Run(){
    last := time.Now()
    for{
        select{
        case response := <- d.Ping:
            response <- "delay"
        default:
        }
        if len(d.in) == 1 && len(d.out) == 0 && last.Add(time.Second * 1).Before(time.Now()) {
            last = time.Now()
            m := <- d.in
            d.out <- m
        }
    }
}

type Logger struct {
    Ping chan chan string
    in chan interface{}
}

func (l *Logger) Run(){
    for{
        select{
        case response := <- l.Ping:
            response <- "logger"
        default:
        }
        if len(l.in) == 1 {
            m := <- l.in
            fmt.Println(m)
        }
    }
} 

func main(){
    pusherToDelay := make(chan interface{}, 1)
    delayToLog := make(chan interface{}, 1)

    pusherPing := make(chan chan string, 1)
    delayPing := make(chan chan string, 1)
    loggerPing := make(chan chan string, 1)

    p := Pusher{
         out: pusherToDelay,
         Ping: pusherPing,
    }

    d := Delay{
        in: pusherToDelay,
        Ping: delayPing,
        out: delayToLog,
    }

    l := Logger{
        in: delayToLog,
        Ping: loggerPing,
    }

    go p.Run()
    go l.Run()
    go d.Run()
    t := time.NewTicker(500 * time.Millisecond)
    for{
        <-t.C
        pushResp := make(chan string)
        delayResp := make(chan string)
        loggerResp := make(chan string)
        pusherPing <- pushResp 
        delayPing <- delayResp
        loggerPing <-loggerResp
        m1 := <-pushResp
        m2 := <-delayResp
        m3 := <-loggerResp
        fmt.Println(m1, m2, m3)
    }
}
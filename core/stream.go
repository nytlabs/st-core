package core

import (
	"log"
	"strconv"
	"sync"

	"github.com/bitly/go-nsq"
)

type Stream struct {
	quit        chan bool
	Out         chan Message // this channel is used by any block that would like to receive messages
	topic       string
	channel     string
	lookupdAddr string
	maxInFlight string
	sync.Mutex
}

func (s Stream) GetType() SourceType {
	return STREAM
}

func (s *Stream) SetSourceParameter(name, value string) {
	switch name {
	case "topic":
		s.topic = value
		log.Println("set stream topic")
	case "channel":
		s.channel = value
		log.Println("set stream channel")
	case "lookupdAddr":
		s.lookupdAddr = value
		log.Println("set stream lookupdAddr")
	case "maxInFlight":
		s.maxInFlight = value
	}
}

func (s *Stream) Describe() map[string]string {
	return map[string]string{
		"topic":       s.topic,
		"channel":     s.channel,
		"lookupAddr":  s.lookupdAddr,
		"maxInFlight": s.maxInFlight,
	}
}

func NewStream() *Stream {
	out := make(chan Message)
	stream := Stream{
		quit:        make(chan bool),
		Out:         out,
		maxInFlight: "10",
	}
	return &stream
}

func (s Stream) Serve() {
	conf := nsq.NewConfig()
	m, err := strconv.Atoi(s.maxInFlight)
	if err != nil {
		log.Println(err)
	} else {
		conf.MaxInFlight = m
	}
	running := false
	reader, err := nsq.NewConsumer(s.topic, s.channel, conf)
	if err != nil {
		log.Println(err)
		log.Println("NSQ Reader is waiting for restart")
		goto Wait
	}

	reader.AddHandler(s)
	err = reader.ConnectToNSQLookupd(s.lookupdAddr)
	if err != nil {
		log.Println(err)
		log.Println("NSQ Reader is waiting for restart")
		goto Wait
	}

	running = true

	// if the reader fails for whatever reason, we need to wait for the user
	// to update the NSQ params.
Wait:
	<-s.quit // this blocks until the stream Source is stopped
	if running {
		reader.Stop()
		<-reader.StopChan // this blocks until the reader is definitely dead
	}
}

func (s Stream) HandleMessage(message *nsq.Message) error {
	s.Out <- string(message.Body)
	return nil
}

func (s Stream) Stop() {
	s.quit <- true
}

// StreamRecieve receives messages from the Stream store.
//
// OutPin 0: received message
func StreamReceive() Spec {
	return Spec{
		Name: "streamReceive",
		Outputs: []Pin{
			Pin{"out"},
		},
		Source: STREAM,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			stream := s.(*Stream)
			select {
			case out[0] = <-stream.Out:
			case f := <-i:
				return f
			}
			return nil
		},
	}
}

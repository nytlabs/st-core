package core

import (
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/bitly/go-nsq"
)

// Stream represents a two-way communication channel with the world outside of streamtools
type Stream struct {
	sync.Mutex
	Out    chan Message
	In     chan Message
	quit   chan bool
	reader NSQReader
	writer NSQWriter
}

func NewStream() {
	out := make(chan Message)
	in := make(chan Message)
	reader := NewNSQReader(out)
	writer := NewNSQWriter(in)
	return Stream{
		Out:    out,
		In:     in,
		quit:   make(chan bool),
		reader: reader,
		writer: writer,
	}
}

func (s *Stream) Start(supervisor *suture.Supervisor) (readerToken, writerToken suture.Token) {
	readerToken := super.Add(s.reader)
	writerToken := super.Add(s.writer)
	return readerToken, writerToken
}

func (s *Stream) RestartWriter(supervisor *suture.Supervisor, writerToken suture.Token) suture.Token {
	supervisor.
}

type NSQReader struct {
	topic       string
	channel     string
	lookupdAddr string
	quit        chan bool
	out         chan interface{}
}

func NewNSQReader(out chan Message) NSQReader {
	return NSQReader{
		out:  out,
		quit: make(chan bool),
	}
}

// Serve handles reading messages from NSQ
func (r *NSQReader) Serve() {
	var reader *nsq.Consumer
	var err error
	conf := nsq.NewConfig()
	reader, err = nsq.NewConsumer(r.topic, r.channel, conf)
	if err != nil {
		log.Println(err)
	} else {
		reader.AddHandler(r)
		err = reader.ConnectToNSQLookupd(r.lookupdAddr)
		if err != nil {
			log.Println(err)
		}
	}
	<-r.quit
}

func (r *NSQReader) HandleMessage(message *nsq.Message) error {
	r.out <- message.Body
	return nil
}

func (r *NSQReader) Stop() {
	r.quit <- true
}

type NSQWriter struct {
	nsqdTCPAddrs string
	topic        string
	quit         chan bool
}

// ServeNSQWriter handles writing messages to NSQ
func (s *Stream) ServeNSQWriter() {
	var writer *nsq.Producer
	var err error
	conf := nsq.NewConfig()
	writer, err = nsq.NewProducer(s.nsqdTCPAddrs, conf)
	if err != nil {
	}
	for {
		select {
		case msg := <-s.Out:
			msgBytes, err := json.Marshal(msg)
			if err != nil {
				log.Println(err)
				continue
			}
			if len(msgBytes) == 0 {
				continue
			}
			err = writer.Publish(s.topic, msgBytes)
			if err != nil {
			}
		case <-s.quitWriter:
			return
		}
	}
}

func (s *Stream) SetSourceParameter(param, value string) error {
	s.Lock()
	switch param {
	case "channel":
		s.channel = param
	case "topic":
		s.topic = param
	case "lookupdAddr":
		s.lookupdAddr = param
	case "nsqdTCPAddrs":
		s.nsqdTCPAddrs = param
	default:
		return errors.New("unknown parameter")
	}
	s.Unlock()
	// TODO this shouldn't start a non-running writer
	s.quitWriter <- true
	go s.ServeNSQWriter()
	return nil
}

// NewStream returns a new stream Store.
func NewStream() Store {
	stream := &Stream{
		Out:        make(chan interface{}),
		In:         make(chan interface{}),
		quitReader: make(chan bool),
		quit:       make(chan bool),
		quitWriter: make(chan bool),
	}
	return stream
}

// Serve starts the Stream service.
// TODO make Reader and Writer objects that themselves implement suture.Service
func (s *Stream) Serve() {
	go s.ServeNSQReader()
	go s.ServeNSQWriter()
	log.Println("serve is waiting for quit")
	<-s.quit
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
		Shared: STREAM,
		Kernel: func(in, out, internal MessageMap, s Store, i chan Interrupt) Interrupt {
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

// StreamSend publishes the inbound message to the Stream store.
//
// InPin 0: message to send
func StreamSend() Spec {
	return Spec{
		Name: "streamSend",
		Inputs: []Pin{
			Pin{"in"},
		},
		Shared: STREAM,
		Kernel: func(in, out, internal MessageMap, s Store, i chan Interrupt) Interrupt {
			stream := s.(*Stream)
			select {
			case stream.In <- in[0]:
			case f := <-i:
				return f
			}
			return nil
		},
	}
}

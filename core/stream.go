package core

import (
	"log"

	"github.com/bitly/go-nsq"
	"github.com/thejerf/suture"
)

type Stream struct {
	quit       chan bool
	in         chan Message
	out        chan Message
	supervisor *suture.Supervisor
	reader     Reader
	writer     Writer
}

func (s *Stream) Serve() {
	s.supervisor.Add(s.reader)
	s.supervisor.Add(s.writer)
	<-quit
}

func (s *Stream) Stop() {
	s.supervisor.Stop()
	quit <- true
}

func (s *Stream) UpdateWriter() {
	// update writer parmas here
	s.writer.Stop()
}

func (s *Stream) UpdateReader() {
	// update reader parmas here
	s.reader.Stop()
}

func NewStream() *Stream {
	super := suture.NewSimple("stream")
	in := make(chan Message)
	out := make(chan Message)
	return &Stream{
		supervisor: super,
		in:         in,
		out:        out,
		quit:       make(chan bool),
		writer:     NewWriter(in),
		reader:     NewReader(out),
	}
}

type Writer struct {
	in   chan Message
	quit chan bool
}

func (w *Writer) Serve() {
	<-w.quit
}

func (w *Writer) Stop() {
	w.quit <- true
}

func NewWriter(in chan Message) *Writer {
	return Writer{
		in:   in,
		quit: make(chan Message),
	}
}

type Reader struct {
	out         chan Message
	quit        chan bool
	topic       string
	channel     string
	lookupdAddr string
}

func (r *Reader) Serve() {
	conf := nsq.NewConfig()
	reader, err := nsq.NewConsumer(r.topic, r.channel, conf)
	if err != nil {
		log.Println(err)
		log.Println("NSQ Reader is waiting for restart")
		goto Wait
	}

	reader.AddHandler(r)
	err = reader.ConnectToNSQLookupd(r.lookupdAddr)
	if err != nil {
		log.Println(err)
		log.Println("NSQ Reader is waiting for restart")
	}
	// if the reader fails for whatever reason, we need to wait for the user
	// to update the NSQ params.
Wait:
	<-r.quit
}

func (r *NSQReader) HandleMessage(message *nsq.Message) error {
	r.out <- message.Body
	return nil
}

func (r *Reader) Stop() {
	r.quit <- true
}

func NewReader(out chan Message, topic, channel, lookupdAddr string) {
	return Reader{
		out:         out,
		quit:        make(chan bool),
		topic:       topic,
		channel:     channel,
		lookupdAddr: lookupdAddr,
	}
}

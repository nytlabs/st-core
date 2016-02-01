package core

import (
	"errors"
	"log"

	"github.com/bitly/go-nsq"
)

func NSQConsumerInterface() SourceSpec {
	return SourceSpec{
		Name: "NSQConsumer",
		Type: NSQCONSUMER,
		New:  NewNSQConsumer,
	}
}

type NSQConsumer struct {
	connectChan chan NSQConf
	topic       string
	fromNSQ     chan string
	quit        chan chan error
}

type NSQConf struct {
	conf       *nsq.Config
	topic      string
	channel    string
	lookupAddr string
	errChan    chan *stcoreError
}

type NSQMsg struct {
	msg     string
	errChan chan *stcoreError
}

func (s NSQConsumer) GetType() SourceType {
	return NSQCONSUMER
}

func NewNSQConsumer() Source {
	return &NSQConsumer{
		quit:        make(chan chan error),
		connectChan: make(chan NSQConf),
		fromNSQ:     make(chan string),
	}
}

func (s NSQConsumer) Serve() {
	var reader *nsq.Consumer
	var err error
	for {
		select {
		case conf := <-s.connectChan:
			reader, err = nsq.NewConsumer(conf.topic, conf.channel, conf.conf)
			if err != nil {
				select {
				case conf.errChan <- NewError("NSQ failed to create Consumer with error:" + err.Error()):
				default:
				}
				continue
			}
			reader.AddHandler(s)
			err = reader.ConnectToNSQLookupd(conf.lookupAddr)
			if err != nil {
				select {
				case conf.errChan <- NewError("NSQ connect failed with:" + err.Error()):
				default:
				}
				continue
			}
		case c := <-s.quit:
			if reader != nil {
				reader.Stop()
				<-reader.StopChan // this blocks until the reader is definitely dead
				reader = nil
			}
			c <- nil
		}
	}
}

func (s NSQConsumer) HandleMessage(message *nsq.Message) error {
	// this blocks until ReceiveMessage is called
	s.fromNSQ <- string(message.Body)
	return nil
}

func (s NSQConsumer) ReceiveMessage(i chan Interrupt) (string, Interrupt, error) {
	// receives message
	select {
	case msg, ok := <-s.fromNSQ:
		if !ok {
			return "", nil, errors.New("NSQ connection has closed")
		}
		return msg, nil, nil
	case f := <-i:
		return "", f, nil
	}
}

func (s NSQConsumer) Stop() {
	m := make(chan error)
	s.quit <- m
	// block until closed
	err := <-m
	if err != nil {
		log.Fatal(err)
	}
}

func NSQConsumerConnect() Spec {
	return Spec{
		Name:    "NSQConsumerConnect",
		Outputs: []Pin{Pin{"connected", BOOLEAN}},
		Inputs:  []Pin{Pin{"topic", STRING}, Pin{"channel", STRING}, Pin{"lookupAddr", STRING}, Pin{"maxInFlight", NUMBER}},
		Source:  NSQCONSUMER,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {

			topic, ok := in[0].(string)
			if !ok {
				out[0] = NewError("NSQConnect requries string topic")
				return nil
			}
			channel, ok := in[1].(string)
			if !ok {
				out[0] = NewError("NSQConnect requries string channel")
				return nil
			}
			lookupAddr, ok := in[2].(string)
			if !ok {
				out[0] = NewError("NSQConnect requries string lookupAddr")
				return nil
			}
			maxInFlight, ok := in[3].(float64)
			if !ok {
				out[0] = NewError("NSQConnect requries number maxInFlight")
				return nil
			}

			conf := nsq.NewConfig()
			conf.MaxInFlight = int(maxInFlight)

			nsq, ok := s.(*NSQConsumer)
			if !ok {
				log.Fatal("could not assert source is NSQ")
			}

			errChan := make(chan *stcoreError)

			connParams := NSQConf{
				conf:       conf,
				topic:      topic,
				channel:    channel,
				lookupAddr: lookupAddr,
				errChan:    errChan,
			}

			nsq.connectChan <- connParams

			// block on connect
			select {
			case err := <-errChan:
				if err != nil {
					out[0] = err
					return nil
				}
				out[0] = true
				return nil
			case f := <-i:
				return f
			}

		},
	}
}

// NSQRecieve receives messages from the NSQ system.
//
// OutPin 0: received message
func NSQConsumerReceive() Spec {
	return Spec{
		Name: "NSQConsumerReceive",
		Outputs: []Pin{
			Pin{"out", STRING},
		},
		Source: NSQCONSUMER,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			nsq := s.(*NSQConsumer)
			msg, f, err := nsq.ReceiveMessage(i)
			if err != nil {
				out[0] = err
				return nil
			}
			if f != nil {
				return f
			}
			out[0] = string(msg)
			return nil
		},
	}
}

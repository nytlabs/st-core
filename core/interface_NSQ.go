package core

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	"github.com/bitly/go-nsq"
)

func NSQInterface() SourceSpec {
	return SourceSpec{
		Name: "NSQClient",
		Type: NSQCLIENT,
		New:  NewNSQ,
	}
}

type NSQ struct {
	connectChan chan NSQConf
	topic       string
	sendChan    chan NSQMsg
	fromNSQ     chan string
	subscribe   chan chan string
	unsubscribe chan chan string
	subscribers map[chan string]struct{}
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

func (s NSQ) GetType() SourceType {
	return NSQCLIENT
}

func NewNSQ() Source {
	return &NSQ{
		quit:        make(chan chan error),
		connectChan: make(chan NSQConf),
		sendChan:    make(chan NSQMsg),
		fromNSQ:     make(chan string),
	}
}

func (s NSQ) Serve() {
	var reader *nsq.Consumer
	var writer *nsq.Producer
	for {
		select {
		case conf := <-s.connectChan:
			reader, err := nsq.NewConsumer(conf.topic, conf.channel, conf.conf)
			if err != nil {
				select {
				case conf.errChan <- NewError("NSQ failed to create Consumer with error:" + err.Error()):
				default:
				}
			}
			reader.AddHandler(s)
			err = reader.ConnectToNSQLookupd(conf.lookupAddr)
			if err != nil {
				select {
				case conf.errChan <- NewError("NSQ connect failed with:" + err.Error()):
				default:
				}
			}
			prodAddr, err := getRandomNode(conf.lookupAddr)
			if err != nil {
				select {
				case conf.errChan <- NewError("getRandomNode failed with:" + err.Error()):
				default:
				}
			}
			log.Println("using", prodAddr, "to publish to")
			writer, err = nsq.NewProducer(prodAddr, conf.conf)
			if err != nil {
				select {
				case conf.errChan <- NewError("creating a new producer failed with:" + err.Error()):
				default:
				}
			}
			s.topic = conf.topic
		case msg := <-s.sendChan:
			if writer == nil {
				msg.errChan <- NewError("NSQ is not connected; cannot send.")
				continue
			}
			err := writer.Publish(s.topic, []byte(msg.msg))
			if err != nil {
				msg.errChan <- NewError("NSQ publish failed with: " + err.Error())
			} else {
				msg.errChan <- nil
			}
		case <-s.quit:
			reader.Stop()
			<-reader.StopChan // this blocks until the reader is definitely dead
			// TODO have some sort of timeout here and return with error maybe?
			// don't forget the object coming through s.quite is an option error channel
			writer.Stop()
		}
	}
}

func (s NSQ) HandleMessage(message *nsq.Message) error {
	// this blocks until ReceiveMessage is called
	s.fromNSQ <- string(message.Body)
	return nil
}

func (s NSQ) ReceiveMessage(i chan Interrupt) (string, Interrupt, error) {
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

func (s NSQ) SendMessage(msg string) *stcoreError {
	// send message
	m := NSQMsg{msg, make(chan *stcoreError)}
	go func() { s.sendChan <- m }()
	err := <-m.errChan
	return err
}

func (s NSQ) Stop() {
	m := make(chan error)
	s.quit <- m
	// block until closed
	err := <-m
	if err != nil {
		log.Fatal(err)
	}
}

type nodesResponse struct {
	Status_code int
	Status_txt  string
	Data        producers
}
type producers struct {
	Producers []producer
}
type producer struct {
	Tcp_port          int
	Broadcast_address string
}

func getRandomNode(lookupdAddr string) (string, error) {
	resp, err := http.Get("http://" + lookupdAddr + "/nodes")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var n nodesResponse
	err = json.Unmarshal(body, &n)
	if err != nil {
		return "", err
	}
	if n.Status_code != 200 {
		return "", errors.New("could not get list of nsqd nodes")
	}
	nProducers := len(n.Data.Producers)
	if nProducers <= 0 {
		log.Fatal(errors.New("found no NSQ daemons"))
	}
	return n.Data.Producers[rand.Intn(nProducers)].Broadcast_address, nil
}

func NSQConnect() Spec {
	return Spec{
		Name:    "NSQConnect",
		Outputs: []Pin{Pin{"connected", BOOLEAN}},
		Inputs:  []Pin{Pin{"topic", STRING}, Pin{"channel", STRING}, Pin{"lookupAddr", STRING}, Pin{"maxInFlight", NUMBER}},
		Source:  NSQCLIENT,
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
				out[0] = NewError("NSQConnect requries number lookupAddr")
				return nil
			}

			conf := nsq.NewConfig()
			conf.MaxInFlight = int(maxInFlight)

			nsq, ok := s.(*NSQ)
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
func NSQReceive() Spec {
	return Spec{
		Name: "NSQReceive",
		Outputs: []Pin{
			Pin{"out", STRING},
		},
		Source: NSQCLIENT,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			nsq := s.(*NSQ)
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

func NSQSend() Spec {
	return Spec{
		Name: "NSQSend",
		Inputs: []Pin{
			Pin{"msg", STRING},
		},
		Outputs: []Pin{
			Pin{"sent", BOOLEAN},
		},
		Source: NSQCLIENT,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			nsq := s.(*NSQ)

			msg, ok := in[0].(string)
			if !ok {
				out[0] = NewError("NSQSend requires string msg")
				return nil
			}

			err := nsq.SendMessage(msg)
			if err != nil {
				out[0] = err
				return nil
			}

			out[0] = true
			return nil
		},
	}
}

package core

import (
	"errors"
	"io"
	"log"
	"strings"

	"golang.org/x/net/websocket"
)

func WebsocketClient() SourceSpec {
	return SourceSpec{
		Name: "wsClient",
		Type: WSCLIENT,
		New:  NewWsClient,
	}
}

const (
	DISCONNECTED = iota
	CONNECTED
)

func (ws *wsClient) SetSourceParameter(name, value string) {
}

func (ws *wsClient) Describe() []map[string]string {
	return []map[string]string{}
}

func (ws wsClient) GetType() SourceType {
	return WSCLIENT
}

type connParams struct {
	url     string
	origin  string
	errChan chan *stcoreError
}

type wsMsg struct {
	msg     string
	errChan chan *stcoreError
}

type wsClient struct {
	IsRunning     bool
	conn          *websocket.Conn
	connectChan   chan connParams
	subscribe     chan chan string
	unsubscribe   chan chan string
	subscribers   map[chan string]struct{}
	quit          chan chan error
	stopReader    chan bool
	fromWebsocket chan string
	sendChan      chan wsMsg
}

func NewWsClient() Source {
	ws := &wsClient{
		IsRunning:     false,
		quit:          make(chan chan error),
		stopReader:    make(chan bool),
		subscribe:     make(chan chan string),
		unsubscribe:   make(chan chan string),
		subscribers:   make(map[chan string]struct{}),
		fromWebsocket: make(chan string),
		sendChan:      make(chan wsMsg),
		connectChan:   make(chan connParams),
	}
	return ws
}

func (ws *wsClient) Serve() {
	var err error
	ws.IsRunning = true
	for {
		select {
		case p := <-ws.connectChan:
			if ws.conn != nil {
				ws.stopReader <- true
			}
			err = ws.Connect(p)

			// these sends need to be non-blocking in case the connect block has been interrupted
			if err != nil {
				select {
				case p.errChan <- NewError("websocket connect failed with:" + err.Error()):
				default:
				}
				break
			}
			select {
			case p.errChan <- nil:
			default:
			}
			go ws.ReadLoop()
		case msg := <-ws.sendChan:
			if ws.conn == nil {
				msg.errChan <- NewError("websocket connection is nil, cannot send")
				continue
			}
			err := websocket.Message.Send(ws.conn, msg.msg)
			if err != nil {
				msg.errChan <- NewError("websocket send failed with: " + err.Error())
			} else {
				msg.errChan <- nil
			}
		case msg := <-ws.fromWebsocket:
			for c, _ := range ws.subscribers {
				c <- msg
			}
		case c := <-ws.subscribe:
			ws.subscribers[c] = struct{}{}
		case c := <-ws.unsubscribe:
			delete(ws.subscribers, c)
		case r := <-ws.quit:
			var err error
			if ws.conn != nil {
				ws.stopReader <- true
			}
			for c, _ := range ws.subscribers {
				delete(ws.subscribers, c)
				close(c)
			}
			if ws.conn != nil {
				err = ws.conn.Close()
			}
			r <- err
			return
		}
	}
}

func (ws wsClient) ReadLoop() {
	c := make(chan string)
	go func() {
		for {
			if ws.conn == nil {
				return
			}
			var msg string
			err := websocket.Message.Receive(ws.conn, &msg)
			if err != nil {
				if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection") {
					m := make(chan error)
					ws.quit <- m
					<-m
					return
				}
				continue
			}
			c <- string(msg)
		}
	}()
	for {
		select {
		case <-ws.stopReader:
			return
		case msg := <-c:
			ws.fromWebsocket <- msg
		}
	}
}

func (ws *wsClient) Stop() {
	ws.IsRunning = false
	m := make(chan error)
	ws.quit <- m
	// block until closed
	err := <-m
	if err != nil {
		log.Fatal(err)
	}
}

func (ws wsClient) ReceiveMessage(i chan Interrupt) (string, Interrupt, error) {
	// receives message
	c := make(chan string, 10)
	ws.subscribe <- c
	select {
	case msg, ok := <-c:
		if !ok {
			return "", nil, errors.New("websocket connection has closed")
		}
		ws.unsubscribe <- c
		return msg, nil, nil
	case f := <-i:
		ws.unsubscribe <- c
		return "", f, nil
	}
}

func (ws wsClient) SendMessage(msg string) *stcoreError {
	if !ws.IsRunning {
		return NewError("cannot send on stopped websocketClient")
	}
	// send message
	m := wsMsg{msg, make(chan *stcoreError)}
	go func() { ws.sendChan <- m }()
	err := <-m.errChan
	return err
}

func (ws *wsClient) Connect(p connParams) error {
	// dial the websocket
	conn, err := websocket.Dial(p.url, "", p.origin)
	if err != nil {
		return err
	}
	ws.conn = conn
	return nil
}

func wsClientConnect() Spec {
	return Spec{
		Name:    "wsClientConnect",
		Outputs: []Pin{Pin{"connected", BOOLEAN}},
		Inputs:  []Pin{Pin{"url", STRING}, Pin{"origin", STRING}},
		Source:  WSCLIENT,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {

			url, ok := in[0].(string)
			if !ok {
				out[0] = NewError("wsClientConnect requries string url")
				return nil
			}

			origin, ok := in[1].(string)
			if !ok {
				out[0] = NewError("wsClientConnect requries string url")
				return nil
			}

			ws := s.(*wsClient)

			errChan := make(chan *stcoreError)
			ws.connectChan <- connParams{url, origin, errChan}

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

func wsClientReceive() Spec {
	return Spec{
		Name:    "wsClientReceive",
		Outputs: []Pin{Pin{"msg", STRING}},
		Source:  WSCLIENT,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			ws := s.(*wsClient)
			msg, f, err := ws.ReceiveMessage(i)
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

func wsClientSend() Spec {
	return Spec{
		Name:    "wsClientSend",
		Inputs:  []Pin{Pin{"msg", STRING}},
		Outputs: []Pin{Pin{"sent", BOOLEAN}},
		Source:  WSCLIENT,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			ws := s.(*wsClient)

			msg, ok := in[0].(string)
			if !ok {
				out[0] = NewError("wsClientSend requires string msg")
				return nil
			}

			err := ws.SendMessage(msg)
			if err != nil {
				out[0] = err
				return nil
			}

			out[0] = true
			return nil

		},
	}
}

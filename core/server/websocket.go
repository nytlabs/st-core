package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type socket struct {
	ws   *websocket.Conn
	send chan []byte
}

func (c *socket) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (s *Server) websocketBroadcast(v interface{}) {
	out, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	s.broadcast <- out
}

func (s *Server) websocketRouter() {
	hub := make(map[*socket]bool)
	for {
		select {
		case c := <-s.addSocket:
			hub[c] = true
		case c := <-s.delSocket:
			delete(hub, c)
		case m := <-s.broadcast:
			for c := range hub {
				c.send <- m
			}
		}
	}
}
func (s *Server) websocketReadPump(c *socket) {
	defer func() {
		s.delSocket <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if string(message) == "list" {
			s.Lock()
			blocks, _ := json.Marshal(s.ListBlocks())
			connections, _ := json.Marshal(s.ListConnections())
			s.Unlock()
			c.send <- blocks
			c.send <- connections
		}

		if err != nil {
			break
		}
		s.emitChan <- message
	}
}

func (s *Server) websocketWritePump(c *socket) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (s *Server) UpdateSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &socket{send: make(chan []byte, 256), ws: ws}
	s.addSocket <- c
	go s.websocketWritePump(c)
	go s.websocketReadPump(c)
}

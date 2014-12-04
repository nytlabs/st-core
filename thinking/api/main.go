package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/thejerf/suture"
)

type Block struct {
	q     chan bool
	token suture.ServiceToken
	id    int
}

func (b *Block) Serve() {
	<-b.q
}

func (b *Block) Stop() {
	b.q <- true
}

func IDGen() func() int {
	i := 0
	return func() int {
		i += 1
		return i
	}
}

func NewBlock(id int) *Block {
	q := make(chan bool)
	return &Block{
		q: q,
	}
}

type Server struct {
	blocks     map[int]*Block
	supervisor *suture.Supervisor
	sync.Mutex
}

func NewServer() *Server {
	supervisor := suture.NewSimple("st-core")
	supervisor.ServeBackground()
	blocks := make(map[int]*Block)
	return &Server{
		supervisor: supervisor,
		blocks:     blocks,
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("hi!"))
}

func (s *Server) createBlockHandler(w http.ResponseWriter, r *http.Request) {
	nextID := IDGen()
	b := NewBlock(nextID())

	b.token = s.supervisor.Add(b)

	s.Lock()
	s.blocks[b.id] = b
	s.Unlock()

	w.WriteHeader(200)
	w.Write([]byte("OK"))

}

func (s *Server) deleteBlockHandler(w http.ResponseWriter, r *http.Request) {
	var blockid int
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, &blockid)
	if err != nil {
		log.Fatal(err)
	}
	block := s.blocks[blockid]

	s.supervisor.Remove(block.token)
	s.Lock()
	delete(s.blocks, blockid)
	s.Unlock()

	w.WriteHeader(200)
	w.Write([]byte("OK"))

}

func (s *Server) blockInfoHandler(w http.ResponseWriter, r *http.Request) {
	var blockid int
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, &blockid)
	if err != nil {
		log.Fatal(err)
	}
	block := s.blocks[blockid]
	w.WriteHeader(200)
	out, err := json.Marshal(block)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(out)
}

func main() {

	s := NewServer()

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/new", s.createSessionHandler).Methods("POST")
	r.HandleFunc("/{session}/blocks", s.createBlockHandler).Methods("POST")
	r.HandleFunc("/{session}blocks/{id}", s.blockInfoHandler).Methods("GET")
	r.HandleFunc("/{session}blocks/{id}", s.deleteBlockHandler).Methods("DELETE")

	http.Handle("/", r)

	log.Println("serving on 7071")

	err := http.ListenAndServe(":7071", nil)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

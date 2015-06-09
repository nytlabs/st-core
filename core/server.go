package core

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

func ServerSource() SourceSpec {
	return SourceSpec{
		Name: "server",
		Type: SERVER,
		New:  NewServer,
	}
}

type handlerRegistration struct {
	c    chan Request
	name string
}

type Server struct {
	quit          chan bool
	router        *mux.Router
	routes        map[string]chan Request
	addHandler    chan handlerRegistration
	removeHandler chan string
	port          int
	sync.Mutex
}

func (s Server) GetType() SourceType {
	return SERVER
}

func (s *Server) SetSourceParameter(name, value string) {
	switch name {
	}
}

func (s *Server) Describe() map[string]string {
	return map[string]string{"port": strconv.Itoa(s.port)}
}

func NewServer() Source {
	server := &Server{
		quit:          make(chan bool),
		router:        mux.NewRouter().StrictSlash(true),
		routes:        make(map[string]chan Request),
		addHandler:    make(chan handlerRegistration),
		removeHandler: make(chan string),
	}

	return server
}

func (s *Server) Serve() {

	server := &http.Server{
		Addr:           ":8080",
		Handler:        s.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// base router
	s.router.HandleFunc("/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		outChan, ok := s.routes[name]
		if !ok {
			log.Println("404")
			return
		}
		outChan <- Request{w, r}
	}).Methods("GET")

	log.Println("starting HTTP server on", server.Addr)
	go server.ListenAndServe()
	for {
		select {
		case <-s.quit:
			return
		case f := <-s.addHandler:
			log.Println("registring new handler", f.name)
			s.routes[f.name] = f.c
		case name := <-s.removeHandler:
			log.Println("removing handler", name)
			delete(s.routes, name)
		}
	}
}

func (s Server) Stop() {
	s.quit <- true
}

type Request struct {
	responseWriter http.ResponseWriter
	request        *http.Request
}

// OutPin 0: received request
func FromRequest() Spec {
	return Spec{
		Name: "endpoint",
		Inputs: []Pin{
			Pin{"name", STRING},
		},
		Outputs: []Pin{
			Pin{"request", OBJECT},
			Pin{"writer", WRITER},
		},
		Source: SERVER,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			log.Println("running server Kernel")
			server := s.(*Server)
			name, ok := in[0].(string)
			if !ok {
				log.Fatal("inputs must be strings")
			}

			requests := make(chan Request)

			log.Println("trying to register")
			server.addHandler <- handlerRegistration{requests, name}
			log.Println("waiting for something")
			select {
			case r := <-requests:
				out[0] = r.request
				out[1] = r.responseWriter
			case f := <-i:
				log.Println("INTTERRRRRUPT")
				server.removeHandler <- name
				return f
			}
			log.Println("done")
			return nil
		},
	}
}

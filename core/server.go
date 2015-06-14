package core

import (
	"io/ioutil"
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

type Request struct {
	responseWriter http.ResponseWriter
	respChan       chan error
	request        *http.Request
}

func (r Request) Write(p []byte) (n int, err error) {
	nn, err := r.responseWriter.Write(p)
	return nn, err
}

func (r Request) Close() error {
	r.respChan <- nil
	return nil
}

func (r Request) Flush() {
	if flusher, ok := r.responseWriter.(http.Flusher); ok {
		flusher.Flush()
	} else {
		log.Println("responseWriter can't flush")
	}
}

func (s *Server) Serve() {

	server := &http.Server{
		Addr:        ":8080",
		Handler:     s.router,
		ReadTimeout: 10 * time.Second,
		//WriteTimeout:   0 * time.Second,
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
		c := make(chan error)
		outChan <- Request{w, c, r}
		err := <-c
		if err != nil {
			log.Println(err)
		}
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
			Pin{"body", STRING},
		},
		Source: SERVER,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			server := s.(*Server)
			name, ok := in[0].(string)
			if !ok {
				log.Fatal("inputs must be strings")
			}

			requests := make(chan Request)

			server.addHandler <- handlerRegistration{requests, name}
			select {
			case r := <-requests:
				body, err := ioutil.ReadAll(r.request.Body)
				if err != nil {
					out[0] = NewError("could not read body")
					return nil
				}
				out[0] = r.request
				out[1] = r
				out[2] = string(body)
			case f := <-i:
				server.removeHandler <- name
				return f
			}
			return nil
		},
	}
}

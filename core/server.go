package core

import (
	"log"
	"net/http"
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

type Server struct {
	quit     chan bool
	requests chan Request // contains the request and the response writer
	router   *mux.Router
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
	return map[string]string{}
}

func NewServer() Source {
	server := &Server{
		quit:     make(chan bool),
		requests: make(chan Request),
		router:   mux.NewRouter().StrictSlash(true),
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
	log.Println("starting listenAndServe")
	go func() {
		log.Fatal(server.ListenAndServe())
	}()
	log.Println("ListenAndServe started")
	<-s.quit
}

func (s Server) Stop() {
	s.quit <- true
}

// registers a new endpoint, or updates an existing endpoint
func (s Server) RegisterEndpoint(name, method, endpoint string) {
	log.Println("Registering endpoint")
	var route *mux.Route
	log.Println("getting route", name)
	route = s.router.Get(name)
	if route == nil {
		log.Println("server source: making new route with name", name)
		route = s.router.NewRoute().Name(name)
	}
	log.Println("populating route")
	route.Methods(method).Path(endpoint).HandlerFunc(s.handler)
	log.Println("route creation complete")
}

type Request struct {
	responseWriter http.ResponseWriter
	request        *http.Request
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	s.requests <- Request{w, r}
}

// OutPin 0: received request
func FromRequest() Spec {
	return Spec{
		Name: "FromRequest",
		Inputs: []Pin{
			Pin{"name", STRING},
			Pin{"method", STRING},
			Pin{"endpoint", STRING},
		},
		Outputs: []Pin{
			Pin{"request", OBJECT},
		},
		Source: SERVER,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			log.Println("running server Kernel")
			server := s.(*Server)
			name, ok := in[0].(string)
			method, ok := in[1].(string)
			endpoint, ok := in[2].(string)
			if !ok {
				log.Fatal("inputs must be strings")
			}
			log.Println("requesting endpoint registration for", name)
			server.RegisterEndpoint(name, method, endpoint)
			log.Println("endpoint", name, "registered")
			select {
			case out[0] = <-server.requests:
			case f := <-i:
				return f
			}
			return nil
		},
	}
}

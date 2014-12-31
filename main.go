package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nytlabs/st-core/core"
)

func main() {

	s := core.NewServer()
	r := mux.NewRouter()
	r.HandleFunc("/", s.WebsocketHandler).Methods("GET")
	r.HandleFunc("/group", s.GetGroupHandler).Methods("GET")
	r.HandleFunc("/group", s.CreateGroupHandler).Methods("POST")
	r.HandleFunc("/block", s.CreateBlockHandler).Methods("POST")
	r.HandleFunc("/connection", s.CreateConnectionHandler).Methods("POST")
	http.Handle("/", r)

	log.Println("serving on 7071")
	err := http.ListenAndServe(":7071", nil)
	if err != nil {
		log.Panicf(err.Error())
	}
}

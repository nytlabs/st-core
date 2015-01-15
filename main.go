package main

import (
	"log"
	"net/http"

	"github.com/nytlabs/st-core/core/server"
)

func main() {

	s := server.NewServer()
	r := s.NewRouter()

	http.Handle("/", r)

	log.Println("serving on 7071")
	err := http.ListenAndServe(":7071", nil)
	if err != nil {
		log.Panicf(err.Error())
	}
}

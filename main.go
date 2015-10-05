package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/nytlabs/st-core/server"
)

func main() {

	// Unpack settings file, or create a new one if necessary
	var settings server.Settings

	dir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	fname := dir + "/.st.json"
	d, err := ioutil.ReadFile(fname)
	if err != nil {
		if os.IsNotExist(err) {
			// make a fresh settings file
			log.Println("creating new settings file at", fname)
			newSettings := server.NewSettings()
			d, err = json.Marshal(newSettings)
			err = ioutil.WriteFile(fname, d, 0644)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}
	err = json.Unmarshal(d, &settings)
	if err != nil {
		log.Fatal(err)
	}

	s := server.NewServer(settings)
	r := s.NewRouter()

	http.Handle("/", r)

	log.Println("serving on 7071")
	err = http.ListenAndServe(":7071", nil)
	if err != nil {
		log.Panicf(err.Error())
	}
}

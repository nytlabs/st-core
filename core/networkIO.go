package core

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

// HTTPGET makes an HTTP GET request to the specified URL, emitting the response as bytes. You will probably want to decode this using JSONUnmarshal or XMLUnmarshal or StringUnmarshal as required.
//
//Pin 0: URL string
//
//Pin 1: header []string
func HTTPGET() Spec {
	return Spec{
		Inputs:  []Pin{Pin{"URL"}, Pin{"header"}},
		Outputs: []Pin{Pin{"response"}},
		Kernel: func(in, out, internal MessageMap, s Store, i chan Interrupt) Interrupt {

			url, ok := in[0].(string)
			if !ok {
				out[0] = NewError("HTTPGET requires url to be a string")
				return nil
			}

			// header should be provided as a map like {"Content-Type": "application/x-www-form-urlencoded"}
			// TODO
			header, ok := in[1].(map[string]string)
			if !ok {
				out[0] = NewError("HTTPGET requres headers to be a map")
				return nil
			}

			// let's only make one client. We'll store it in the internal state
			var client http.Client
			clientInterface, ok := internal[0]
			if ok {
				client, ok = clientInterface.(http.Client)
				if !ok {
					log.Fatal("found non-client in the internal state")
				}
			} else {
				// TODO do we want to specify timeout as a Pin here? Would we have to ditch the idea of reusing the client?
				transport := http.Transport{
					Dial: func(network, addr string) (net.Conn, error) {
						return net.DialTimeout(network, addr, time.Duration(10*time.Second))
					},
				}
				client = http.Client{
					Transport: &transport,
				}
				internal[0] = client
			}

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Fatal(err)
			}
			for key, value := range header {
				if key == "Host" {
					req.Host = value
				} else {
					req.Header.Set(key, value)
				}
			}

			resp, err := client.Do(req)
			if err != nil {
				out[0] = err
				return nil
			}

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				out[0] = err
				return nil
			}
			out[0] = body

			return nil
		},
	}
}

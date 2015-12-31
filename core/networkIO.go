package core

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// HTTPResponse makes an HTTP request to the specified URL, emitting the response as a string.
func HTTPRequest() Spec {
	return Spec{
		Name:    "HTTPRequest",
		Inputs:  []Pin{Pin{"URL", STRING}, Pin{"header", OBJECT}, Pin{"method", STRING}, Pin{"body", STRING}},
		Outputs: []Pin{Pin{"response", STRING}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {

			url, ok := in[0].(string)
			if !ok {
				out[0] = NewError("HTTPResponse requires url to be a string")
				return nil
			}

			// header should be provided as a map like {"Content-Type": "application/x-www-form-urlencoded"}
			// TODO
			header, ok := in[1].(map[string]interface{})
			if !ok {
				out[0] = NewError("HTTPResponse requres headers to be an object")
				return nil
			}
			method, ok := in[2].(string)
			if !ok {
				out[0] = NewError("HTTPRequest requires method to be a string")
				return nil
			}

			ok = false
			for _, m := range []string{"GET", "POST", "PUT", "DELETE", "HEAD", "TRACE", "OPTIONS", "PATCH"} {
				if m == method {
					ok = true
					break
				}
			}
			if !ok {
				out[0] = NewError("HTTPRequest does not support requested Method")
				return nil
			}

			requestBody, ok := in[3].(string)
			if !ok {
				out[0] = NewError("HTTPRequest requires a string body")
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

			req, err := http.NewRequest(method, url, strings.NewReader(requestBody))
			if err != nil {
				out[0] = NewError("Could not build request")
				return nil
			}
			for key, value := range header {
				vstring, ok := value.(string)
				if !ok {
					out[0] = NewError("Header values must be strings")
					return nil
				}
				if key == "Host" {
					req.Host = vstring
				} else {
					req.Header.Set(key, vstring)
				}
			}
			errChan := make(chan error)
			respChan := make(chan *http.Response)
			go func() {
				resp, err := client.Do(req)
				if err != nil {
					errChan <- err
				}
				respChan <- resp
			}()

			select {
			case err := <-errChan:
				out[0] = err
				return nil
			case resp := <-respChan:
				defer resp.Body.Close()
				responseBody, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					out[0] = err
					return nil
				}
				out[0] = string(responseBody)
				return nil
			case f := <-i:
				t := client.Transport
				transport := t.(*http.Transport)
				transport.CancelRequest(req)
				return f

			}

			return nil

		},
	}
}

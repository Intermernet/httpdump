package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var port = ":8080"

type server struct{}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request received from %s\n", r.RemoteAddr)
	log.Println("Headers:")
	for k, v := range r.Header {
		fmt.Printf("%s\t%s\n", k, v)
	}
	if r.ContentLength > 0 {
		log.Println("Body:")
		defer r.Body.Close()
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("%v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		switch {
		case strings.HasSuffix(r.Header["Content-Type"][0], "json"):
			var out bytes.Buffer
			if err := json.Indent(&out, b, "", "  "); err != nil {
				log.Printf("%v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			fmt.Printf("%s\n", out.Bytes())
		default:
			log.Printf("%s\n", b)
		}
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	log.Printf("Listening on port %s\n", port)
	s := new(server)
	http.Handle("/", s)
	log.Fatal(http.ListenAndServe(port, nil))
}

package main

import (
	"compress/gzip"
	"io/ioutil"
	"log"
	"net/http"
)

var port = ":8080"

type server struct{}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v\n", r.Header)
	defer r.Body.Close()
	for _, enc := range r.Header["Accept-Encoding"] {
		switch {
		case enc == "gzip":
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Printf("%v\n", err)
			}
			b, err := ioutil.ReadAll(gz)
			if err != nil {
				log.Printf("%v\n", err)
			}
			log.Printf("%v\n", b)
		default:
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Printf("%v\n", err)
			}
			log.Printf("%v\n", b)
		}
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	log.Printf("Listening on port %s\n", port)
	s := new(server)
	log.Fatal(http.ListenAndServe(port, s))
}

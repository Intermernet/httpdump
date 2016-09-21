package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

var port = ":8080"

type server struct{}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v\n", r.Header)
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("%v\n", err)
	}
	log.Printf("%s\n", b)
	w.WriteHeader(http.StatusOK)
}

func main() {
	log.Printf("Listening on port %s\n", port)
	s := new(server)
	log.Fatal(http.ListenAndServe(port, s))
}

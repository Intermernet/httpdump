package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var addr string

func init() {
	flag.StringVar(&addr, "addr", ":8080", "address to listen on (:8080, localhost:1234 etc.)")
}

type server struct{}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request received from %s\n", r.RemoteAddr)
	fmt.Println("Headers:")
	for k, v := range r.Header {
		fmt.Printf("%s\t%s\n", k, v)
	}
	if r.ContentLength > 0 {
		fmt.Println("Body:")
		defer r.Body.Close()
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("%v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Parse Content-Types to check if JSON or variant.
		// See https://www.iana.org/assignments/media-types/media-types.xhtml for full list
		js := false
		for _, cts := range r.Header["Content-Type"] {
			ct := strings.Split(cts, ";")[0]   // Split at ";" to discard charset if defined
			if strings.HasSuffix(ct, "json") { // All JSON content-types should end in "json"
				js = true
			}
		}
		switch {
		case js:
			var out bytes.Buffer
			if err := json.Indent(&out, b, "", "  "); err != nil {
				log.Printf("%v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Printf("%s\n", out.Bytes())
		default:
			log.Printf("%s\n", b)
		}
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	flag.Parse()
	log.Printf("Listening at address %s\n", addr)
	s := new(server)
	http.Handle("/", s)
	log.Fatal(http.ListenAndServe(addr, nil))
}

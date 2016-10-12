package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	addr       string
	reqtimeout time.Duration
)

func init() {
	flag.StringVar(&addr, "addr", ":8080", "address to listen on (:8080, localhost:1234 etc.)")
	flag.DurationVar(&reqtimeout, "timeout", 5*time.Second, "request timeout (5s, 100ms, 1h etc.)")
}

type server struct{}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), reqtimeout)
	req := r.WithContext(ctx)
	defer cancel()
	defer req.Body.Close()
	select {
	default:
		log.Printf("Request received from %s\n", req.RemoteAddr)
		fmt.Println("Headers:")
		for k, v := range req.Header {
			fmt.Printf("%s\t%s\n", k, v)
		}
		if req.ContentLength > 0 {
			fmt.Println("Body:")
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Printf("%v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// Parse Content-Types to check if JSON or variant.
			// See https://www.iana.org/assignments/media-types/media-types.xhtml for full list
			js := false
			for _, cts := range req.Header["Content-Type"] {
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
	case <-ctx.Done():
		log.Printf("Request timeout: %s %s\n", reqtimeout, ctx.Err())
		w.WriteHeader(http.StatusRequestTimeout)
	}
}

func main() {
	flag.Parse()
	log.Printf("Listening at address %s\n", addr)
	s := new(server)
	http.Handle("/", s)
	log.Fatal(http.ListenAndServe(addr, nil))
}

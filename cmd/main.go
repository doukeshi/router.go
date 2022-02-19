package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/doukeshi/router.go"
)

var greet = func(w http.ResponseWriter, r *http.Request) {
	log.Printf("Hello World! %s", time.Now())
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func wrapperFunc(hf http.HandlerFunc, param string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("%s Pre Wrapper", param)
		hf.ServeHTTP(rw, r)
		log.Printf("%s Post Wrapper", param)
	}
}

func main() {
	r := router.New()

	r.HandleFunc(http.MethodGet, "/", greet)
	r.HandleFunc(http.MethodGet, "/ping", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, "pong")
	})

	v1 := r.GroupFunc("/v1", func(rw http.ResponseWriter, r *http.Request, h http.Handler) {
		defer func() func() {
			log.Printf("v1 Pre Handler")
			return func() {
				log.Printf("v1 Post Handler")
			}
		}()()
		h.ServeHTTP(rw, r)
	})
	v1.HandleFunc(http.MethodGet, "/", greet)
	v1.HandleFunc(http.MethodGet, "/go", wrapperFunc(greet, "hey"))
	v1.HandleFunc(http.MethodGet, "/ping", func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("--------")
		fmt.Fprintf(rw, "pong")
	})

	log.Fatal(http.ListenAndServe(":8000", r))
}

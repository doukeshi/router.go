package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/doukeshi/router.go"
)

var greet = func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
	r := router.New()

	r.HandleFunc(http.MethodGet, "/", greet)
	r.HandleFunc(http.MethodGet, "/ping", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, "pong")
	})

	log.Fatal(http.ListenAndServe(":8000", r))
}

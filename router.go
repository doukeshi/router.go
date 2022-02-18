package router

import (
	"fmt"
	"net/http"
)

const (
	MOD string = "router.go"
)

type router struct {
	table map[string]map[string]http.Handler
}

func New() *router {
	return &router{table: make(map[string]map[string]http.Handler)}
}

func (r *router) Handle(method, pattern string, h http.Handler) {
	if _, ok := r.table[pattern]; !ok {
		r.table[pattern] = make(map[string]http.Handler)
	}
	if _, ok := r.table[pattern][method]; ok {
		panic(fmt.Sprintf("%s: handler already registered on [%s %s]", MOD, method, pattern))
	}
	r.table[pattern][method] = h
}

func (r *router) HandleFunc(method, pattern string, hf http.HandlerFunc) {
	r.Handle(method, pattern, hf)
}

func (r *router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path

	ph, ok := r.table[path]
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	h, ok := ph[method]
	if !ok {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	h.ServeHTTP(res, req)
}

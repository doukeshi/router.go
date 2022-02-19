package router

import (
	"net/http"
	"strings"
)

const (
	MOD string = "router.go"
)

type router struct {
	table *tree
}

func New() *router {
	return &router{table: NewTree()}
}

func (r *router) Handle(method, pattern string, h http.Handler) {
	r.table.i(method, pattern, h)
}

func (r *router) HandleFunc(method, pattern string, hf http.HandlerFunc) {
	r.Handle(method, pattern, hf)
}

func (r *router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path

	h := r.table.lookup(method, path)
	h.ServeHTTP(res, req)
}

// error handlers
var err404 = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
})

type err405 struct {
	allowedMethods []string
}

func (err err405) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Allow", strings.Join(err.allowedMethods, ","))
	res.WriteHeader(http.StatusMethodNotAllowed)
}

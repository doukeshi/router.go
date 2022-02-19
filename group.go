package router

import "net/http"

type groupFunc func(method, pattern string, h http.Handler)

func (g groupFunc) Handle(method, pattern string, h http.Handler) {
	g(method, pattern, h)
}

func (g groupFunc) HandleFunc(method, pattern string, hf http.HandlerFunc) {
	g.Handle(method, pattern, hf)
}

type GroupHandler interface {
	ServeGroup(http.ResponseWriter, *http.Request, http.Handler)
}

type GroupHandlerFunc func(http.ResponseWriter, *http.Request, http.Handler)

func (gf GroupHandlerFunc) ServeGroup(res http.ResponseWriter, req *http.Request, h http.Handler) {
	gf(res, req, h)
}

func (r *router) Group(prefix string, gh GroupHandler) groupFunc {
	return func(method, pattern string, h http.Handler) {
		r.Handle(method, prefix+pattern, http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				gh.ServeGroup(rw, r, h)
			},
		))
	}
}

func (r *router) GroupFunc(prefix string, gf GroupHandlerFunc) groupFunc {
	return r.Group(prefix, gf)
}

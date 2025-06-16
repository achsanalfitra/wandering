package router

import "net/http"

// minimal router that impement ServeHTTP

type routeKey struct {
	m string
	p string
}

type Router struct {
	handlers map[routeKey]http.Handler
}

func NewRouter() *Router {
	return &Router{
		handlers: make(map[routeKey]http.Handler),
	}
}

func (r *Router) Register(m, p string, h http.Handler) {
	// register pattern Method, Path, Handler
	r.handlers[routeKey{m: m, p: p}] = h
}

func (r *Router) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	m := rq.Method
	p := rq.URL.Path

	if _, ok := r.handlers[routeKey{m: m, p: p}]; ok {
		r.handlers[routeKey{m: m, p: p}].ServeHTTP(w, rq)
		return
	}

	http.NotFound(w, rq)
}

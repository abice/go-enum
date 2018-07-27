package httputil

import (
	"errors"
	"net/http"
)

var (
	ErrMissingRequiredArgument = errors.New("missing required argument")
)

type HandlerMatcher interface {
	Match(r *http.Request) int
}

type HandlerGet struct{}
type HandlerHead struct{}
type HandlerPost struct{}
type HandlerPut struct{}
type HandlerPatch struct{}
type HandlerDelete struct{}
type HandlerConnect struct{}
type HandlerOptions struct{}
type HandlerTrace struct{}

func (HandlerGet) Match(r *http.Request) int     { return isMethodNotAllowed(r.Method, "GET") }
func (HandlerHead) Match(r *http.Request) int    { return isMethodNotAllowed(r.Method, "HEAD") }
func (HandlerPost) Match(r *http.Request) int    { return isMethodNotAllowed(r.Method, "POST") }
func (HandlerPut) Match(r *http.Request) int     { return isMethodNotAllowed(r.Method, "PUT") }
func (HandlerPatch) Match(r *http.Request) int   { return isMethodNotAllowed(r.Method, "PATCH") }
func (HandlerDelete) Match(r *http.Request) int  { return isMethodNotAllowed(r.Method, "DELETE") }
func (HandlerConnect) Match(r *http.Request) int { return isMethodNotAllowed(r.Method, "CONNECT") }
func (HandlerOptions) Match(r *http.Request) int { return isMethodNotAllowed(r.Method, "OPTIONS") }
func (HandlerTrace) Match(r *http.Request) int   { return isMethodNotAllowed(r.Method, "TRACE") }

type handler struct {
	method string
	h      http.Handler
}

func NewHandler(method string, h http.Handler) http.Handler {
	return handler{
		method: method,
		h:      h,
	}
}

func NewHandlerFunc(method string, h http.HandlerFunc) http.Handler {
	return NewHandler(method, h)
}

func (h handler) Match(r *http.Request) int                        { return isMethodNotAllowed(r.Method, h.method) }
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) { h.h.ServeHTTP(w, r) }

func isMethodNotAllowed(got, expected string) int {
	if got != expected {
		return http.StatusMethodNotAllowed
	}
	return http.StatusOK
}

type ServeMux struct {
	*http.ServeMux
}

func NewServeMux() *ServeMux {
	mux := new(ServeMux)
	mux.ServeMux = http.NewServeMux()
	return mux
}

func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h, _ := mux.Handler(r)
	if matcher, ok := h.(HandlerMatcher); ok {
		if status := matcher.Match(r); status != http.StatusOK {
			w.WriteHeader(status)
			return
		}
	}
	h.ServeHTTP(w, r)
}

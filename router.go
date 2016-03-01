package router

import (
	"net/http"

	"github.com/mosson/lex"
)

// Handler is wrapper of net/http.Handler, using for wrapping simple func as net/http.Handler
type Handler struct {
	handleFn func(http.ResponseWriter, *http.Request)
}

func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.handleFn(w, r)
}

// Router is httpRouter handling path and invoke handler
type Router struct {
	routes map[*lex.Parser]http.Handler
}

// New returns new Router
func New() *Router {
	return &Router{routes: make(map[*lex.Parser]http.Handler)}
}

// RegisterFn registers path pattern to router as func
func (router *Router) RegisterFn(path string, fn func(http.ResponseWriter, *http.Request)) {
	parser := lex.PathParser(path)
	router.routes[&parser] = &Handler{handleFn: fn}
}

// Register registers path pattern to router as Handler
func (router *Router) Register(path string, handler http.Handler) {
	parser := lex.PathParser(path)
	router.routes[&parser] = handler
}

// Handle handles path, registered handler exists, then invokes handler
// Handle returns bool as existing path
func (router *Router) Handle(w http.ResponseWriter, r *http.Request, path string) bool {
	// TODO implements o(n) searching @ mosson/lex
	for parserPtr, handler := range router.routes {
		parser := *parserPtr
		result := parser(path, 0)
		if result.Success {

			for key, value := range result.Attributes {
				if r.Form[key] == nil {
					r.Form[key] = make([]string, 0)
				}
				r.Form[key] = append(r.Form[key], value)
			}

			handler.ServeHTTP(w, r)
			return true
		}
	}

	return false
}

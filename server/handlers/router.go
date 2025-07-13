package handlers

import (
	"jinovatka/assert"
	"net/http"
)

// Serves as default handler intended to be passed directly to server.
// Holds ServeMux with routes from other handlers, and dispatches requests to it.
type RouterHandler struct {
	// Mux shared among the handlers
	Mux *http.ServeMux
}

func NewRouterHandler(mux *http.ServeMux) *RouterHandler {
	return &RouterHandler{
		Mux: mux,
	}
}

// Add handlers to the router
func (router *RouterHandler) AddHandlers(handlers ...Handler) {
	assert.Must(len(handlers) > 0, "RouterHandler.AddHandlers: at least one handler must be passed in") // Yes we crash. This can only happen if I forgot to pass handlers to the function.
	assert.Must(router.Mux != nil, "RouterHandler.AddHandlers: router.Mux can't be nil; only use routers created by NewRouterHandler function")
	for _, handler := range handlers {
		handler.Routes(router.Mux)
	}
}

func (router *RouterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.Mux.ServeHTTP(w, r)
}

// Handlers implementing this interface can have it's routes added to router.
// Use the http.Hanlder interface in cases when the ServeHTTP method is necessary.
type Handler interface {
	// The Routes method should add routes of a handler to the shared ServeMux.
	// It is expected that handlers will also add routes from their subhandlers and so on.
	// This should result in a tree of handlers. Names of the routes should reflect this.
	Routes(*http.ServeMux)
}

package verigo

import (
	"net/http"

	"golang.org/x/net/context"
)

// Middleware represents a constructor for a piece of middleware.
type Middleware func(ContextHandler) ContextHandler

// Chain acts as a list of Middleware constructors.
type Chain struct {
	m []Middleware
}

// New creates a new chain of middlewares.
func New(m ...Middleware) Chain {
	return Chain{m: m}
}

// Then chains the middleware and returns the final http.Handler
func (c Chain) Then(fn func(context.Context, http.ResponseWriter, *http.Request)) http.Handler {
	var final ContextHandler = ContextHandlerFunc(fn)
	for i := len(c.m) - 1; i >= 0; i-- {
		final = c.m[i](final)
	}
	return &ContextAdapter{
		ctx:     context.Background(),
		handler: final,
	}
}

// ContextHandler is the interface similar to http.Handler but allowing to serve Context as well.
type ContextHandler interface {
	ServeHTTPContext(context.Context, http.ResponseWriter, *http.Request)
}

// ContextHandlerFunc is similar to http.HandlerFunc but allowing to serve Context as well.
type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

// ServeHTTPContext calls f(ctx, w, r).
func (f ContextHandlerFunc) ServeHTTPContext(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	f(ctx, w, r)
}

// ContextAdapter represents the adapter for ContextHandler type.
type ContextAdapter struct {
	ctx     context.Context
	handler ContextHandler
}

func (ca *ContextAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ca.handler.ServeHTTPContext(ca.ctx, w, r)
}

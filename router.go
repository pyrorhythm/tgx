package tgx

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

// Router is a composable route tree backed by telegohandler.HandlerGroup.
//
// Outer middleware (Router.Use) maps to telego group middleware and runs before
// route matching within the group. Inner middleware (per-route) runs after a
// route matches, wrapping only that handler.
type Router struct {
	d     *Dispatcher
	group *th.HandlerGroup
}

// Use registers outer middleware on this router's group.
func (r *Router) Use(mw ...Middleware) *Router {
	for _, m := range mw {
		r.group.Use(outerMiddleware(m))
	}
	return r
}

// Group creates a child router. Predicates on the group apply to all routes in the subtree.
func (r *Router) Group(filters ...Filter) *Router {
	g := r.group.Group(toPredicates(filters...)...)
	return &Router{d: r.d, group: g}
}

func outerMiddleware(m Middleware) th.Handler {
	return func(thCtx *th.Context, update telego.Update) error {
		ctx := wrapCtx(thCtx, update)
		return m(ctx, func() error { return thCtx.Next(update) })
	}
}

func (r *Router) register(thHandler th.Handler, filters []Filter, inner []Middleware) {
	h := thHandler
	if len(inner) > 0 {
		h = func(thCtx *th.Context, update telego.Update) error {
			ctx := wrapCtx(thCtx, update)
			return chainInner(ctx, func() error { return thHandler(thCtx, update) }, inner)
		}
	}
	r.group.Handle(h, toPredicates(filters...)...)
}

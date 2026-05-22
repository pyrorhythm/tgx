package tgx

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"

	"pyrorhythm.dev/fn/opt"
)

// RouteOption configures route registration.
type RouteOption func(*routeConfig)

type routeConfig struct {
	filters []Filter
	inner   []Middleware
}

// WithFilters adds route predicates.
func WithFilters(filters ...Filter) RouteOption {
	return func(c *routeConfig) { c.filters = append(c.filters, filters...) }
}

// WithInner adds per-route inner middleware (runs after match).
func WithInner(mw ...Middleware) RouteOption {
	return func(c *routeConfig) { c.inner = append(c.inner, mw...) }
}

func applyRouteOpts(opts ...RouteOption) routeConfig {
	cfg := routeConfig{}
	for _, o := range opts {
		o(&cfg)
	}
	return cfg
}

func registerHandler[T any](
	r *Router,
	cfg routeConfig,
	narrow func(*telego.Update) opt.Of[T],
	h Handler[T],
) {
	r.register(func(thCtx *th.Context, update telego.Update) error {
		ev := narrow(&update)
		if !ev.Valid() {
			return nil
		}
		return h(wrapCtx(thCtx, update), ev.Val())
	}, cfg.filters, cfg.inner)
}

// On registers a typed handler with custom narrow and filters.
func On[T any](r *Router, narrow func(*telego.Update) opt.Of[T], h Handler[T], opts ...RouteOption) {
	registerHandler(r, applyRouteOpts(opts...), narrow, h)
}

// OnMessage registers a message handler.
func (r *Router) OnMessage(h Handler[telego.Message], opts ...RouteOption) {
	registerHandler(r, applyRouteOpts(opts...), NarrowMessage, h)
}

// OnEditedMessage registers an edited message handler.
func (r *Router) OnEditedMessage(h Handler[telego.Message], opts ...RouteOption) {
	registerHandler(r, applyRouteOpts(opts...), NarrowEditedMessage, h)
}

// OnCallbackQuery registers a callback query handler.
func (r *Router) OnCallbackQuery(h Handler[telego.CallbackQuery], opts ...RouteOption) {
	registerHandler(r, applyRouteOpts(opts...), NarrowCallbackQuery, h)
}

// OnInlineQuery registers an inline query handler.
func (r *Router) OnInlineQuery(h Handler[telego.InlineQuery], opts ...RouteOption) {
	registerHandler(r, applyRouteOpts(opts...), NarrowInlineQuery, h)
}

// OnChosenInlineResult registers a chosen inline result handler.
func (r *Router) OnChosenInlineResult(h Handler[telego.ChosenInlineResult], opts ...RouteOption) {
	registerHandler(r, applyRouteOpts(opts...), NarrowChosenInlineResult, h)
}

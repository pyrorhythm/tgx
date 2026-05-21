package tgx

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
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

func updateFromMessage(thCtx *th.Context, msg telego.Message) telego.Update {
	return telego.Update{UpdateID: thCtx.UpdateID(), Message: &msg}
}

func updateFromEditedMessage(thCtx *th.Context, msg telego.Message) telego.Update {
	return telego.Update{UpdateID: thCtx.UpdateID(), EditedMessage: &msg}
}

func updateFromCallback(thCtx *th.Context, q telego.CallbackQuery) telego.Update {
	return telego.Update{UpdateID: thCtx.UpdateID(), CallbackQuery: &q}
}

func updateFromInline(thCtx *th.Context, q telego.InlineQuery) telego.Update {
	return telego.Update{UpdateID: thCtx.UpdateID(), InlineQuery: &q}
}

func updateFromChosen(thCtx *th.Context, ch telego.ChosenInlineResult) telego.Update {
	return telego.Update{UpdateID: thCtx.UpdateID(), ChosenInlineResult: &ch}
}

// On registers a typed handler with custom narrow and filters.
func On[T any](r *Router, narrow func(*telego.Update) (T, bool), h Handler[T], opts ...RouteOption) {
	cfg := applyRouteOpts(opts...)
	r.register(func(thCtx *th.Context, update telego.Update) error {
		ev, ok := narrow(&update)
		if !ok {
			return nil
		}
		return h(wrapCtx(thCtx, update), ev)
	}, cfg.filters, cfg.inner)
}

// OnMessage registers a message handler.
func (r *Router) OnMessage(h Handler[telego.Message], opts ...RouteOption) {
	cfg := applyRouteOpts(opts...)
	thHandler := func(thCtx *th.Context, update telego.Update) error {
		msg, ok := narrowMessage(&update)
		if !ok {
			return nil
		}
		return h(wrapCtx(thCtx, update), msg)
	}
	if len(cfg.inner) == 0 {
		r.group.HandleMessage(func(thCtx *th.Context, msg telego.Message) error {
			return h(wrapCtx(thCtx, updateFromMessage(thCtx, msg)), msg)
		}, toPredicates(cfg.filters...)...)
		return
	}
	r.register(thHandler, cfg.filters, cfg.inner)
}

// OnEditedMessage registers an edited message handler.
func (r *Router) OnEditedMessage(h Handler[telego.Message], opts ...RouteOption) {
	cfg := applyRouteOpts(opts...)
	r.group.HandleEditedMessage(func(thCtx *th.Context, msg telego.Message) error {
		return h(wrapCtx(thCtx, updateFromEditedMessage(thCtx, msg)), msg)
	}, toPredicates(cfg.filters...)...)
}

// OnCallbackQuery registers a callback query handler.
func (r *Router) OnCallbackQuery(h Handler[telego.CallbackQuery], opts ...RouteOption) {
	cfg := applyRouteOpts(opts...)
	thHandler := func(thCtx *th.Context, update telego.Update) error {
		q, ok := narrowCallbackQuery(&update)
		if !ok {
			return nil
		}
		return h(wrapCtx(thCtx, update), q)
	}
	if len(cfg.inner) == 0 {
		r.group.HandleCallbackQuery(func(thCtx *th.Context, q telego.CallbackQuery) error {
			return h(wrapCtx(thCtx, updateFromCallback(thCtx, q)), q)
		}, toPredicates(cfg.filters...)...)
		return
	}
	r.register(thHandler, cfg.filters, cfg.inner)
}

// OnInlineQuery registers an inline query handler.
func (r *Router) OnInlineQuery(h Handler[telego.InlineQuery], opts ...RouteOption) {
	cfg := applyRouteOpts(opts...)
	thHandler := func(thCtx *th.Context, update telego.Update) error {
		q, ok := narrowInlineQuery(&update)
		if !ok {
			return nil
		}
		return h(wrapCtx(thCtx, update), q)
	}
	if len(cfg.inner) == 0 {
		r.group.HandleInlineQuery(func(thCtx *th.Context, q telego.InlineQuery) error {
			return h(wrapCtx(thCtx, updateFromInline(thCtx, q)), q)
		}, toPredicates(cfg.filters...)...)
		return
	}
	r.register(thHandler, cfg.filters, cfg.inner)
}

// OnChosenInlineResult registers a chosen inline result handler.
func (r *Router) OnChosenInlineResult(h Handler[telego.ChosenInlineResult], opts ...RouteOption) {
	cfg := applyRouteOpts(opts...)
	thHandler := func(thCtx *th.Context, update telego.Update) error {
		ch, ok := narrowChosenInlineResult(&update)
		if !ok {
			return nil
		}
		return h(wrapCtx(thCtx, update), ch)
	}
	if len(cfg.inner) == 0 {
		r.group.HandleChosenInlineResult(func(thCtx *th.Context, ch telego.ChosenInlineResult) error {
			return h(wrapCtx(thCtx, updateFromChosen(thCtx, ch)), ch)
		}, toPredicates(cfg.filters...)...)
		return
	}
	r.register(thHandler, cfg.filters, cfg.inner)
}

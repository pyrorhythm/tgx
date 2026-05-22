package tgx

import (
	"context"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"pyrorhythm.dev/fn/res"
)

// Dispatcher wraps telego bot update handling with tgx routers.
type Dispatcher struct {
	Bot     *telego.Bot
	Updates <-chan telego.Update
	bh      *th.BotHandler
	root    *Router
	onError ErrorHandler
}

// DispatcherOption configures the dispatcher.
type DispatcherOption func(*Dispatcher) error

// WithErrorHandler sets a custom error handler on the underlying BotHandler.
func WithErrorHandler(h ErrorHandler) DispatcherOption {
	return func(d *Dispatcher) error {
		d.onError = h
		return nil
	}
}

// NewDispatcher creates a dispatcher from bot and updates channel.
func NewDispatcher(bot *telego.Bot, updates <-chan telego.Update, opts ...DispatcherOption) res.Of[*Dispatcher] {
	d := &Dispatcher{Bot: bot, Updates: updates}

	var thOpts []th.BotHandlerOption
	for _, opt := range opts {
		if err := opt(d); err != nil {
			return res.Err[*Dispatcher](err)
		}
	}
	if d.onError != nil {
		thOpts = append(thOpts, th.WithErrorHandler(func(ctx *th.Context, update telego.Update, err error) {
			d.onError(wrapCtx(ctx, update), err)
		}))
	}

	bh, err := th.NewBotHandler(bot, updates, thOpts...)
	if err != nil {
		return res.Err[*Dispatcher](err)
	}
	d.bh = bh
	d.root = &Router{d: d, group: bh.BaseGroup()}
	return res.OK(d)
}

// Router returns the root router.
func (d *Dispatcher) Router() *Router { return d.root }

// Start begins processing updates.
func (d *Dispatcher) Start() error { return d.bh.Start() }

// Stop stops processing updates.
func (d *Dispatcher) Stop() error { return d.bh.Stop() }

// StopWithContext stops with cancellation.
func (d *Dispatcher) StopWithContext(ctx context.Context) error {
	return d.bh.StopWithContext(ctx)
}

// BotHandler exposes the underlying telegohandler for advanced use.
func (d *Dispatcher) BotHandler() *th.BotHandler { return d.bh }

package tgx

// Handler is a typed update handler.
type Handler[T any] func(ctx *Ctx, ev T) error

// Middleware runs around handlers. Return nil to continue; skip ctx.Next equivalent by not calling next.
type Middleware func(ctx *Ctx, next func() error) error

// ErrorHandler observes handler errors.
type ErrorHandler func(ctx *Ctx, err error)

func chainInner(ctx *Ctx, next func() error, inner []Middleware) error {
	if len(inner) == 0 {
		return next()
	}
	var run func(int) error
	run = func(i int) error {
		if i >= len(inner) {
			return next()
		}
		return inner[i](ctx, func() error { return run(i + 1) })
	}
	return run(0)
}

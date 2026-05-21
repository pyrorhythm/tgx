// Package di provides typed dependency injection via context values.
package di

import (
	"context"

	"pyrorhythm.dev/fn/opt"
)

type key[T any] struct{}

// Put stores a typed value in ctx.
func Put[T any](ctx context.Context, v T) context.Context {
	return context.WithValue(ctx, key[T]{}, v)
}

// Get returns a typed value from ctx.
func Get[T any](ctx context.Context) opt.Of[T] {
	v, ok := ctx.Value(key[T]{}).(T)
	if !ok {
		return opt.Nil[T]()
	}
	return opt.SomeAny(v)
}

package tgx

import (
	"context"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"

	"pyrorhythm.dev/fn/opt"
)

// Filter matches updates. Filters are plain functions (no reflection).
type Filter func(u *telego.Update) bool

// FromOpt converts an optional predicate into a Filter.
func FromOpt(pred func(u *telego.Update) opt.Of[struct{}]) Filter {
	return func(u *telego.Update) bool {
		return pred(u).Valid()
	}
}

// Predicate converts a Filter to a telegohandler.Predicate.
func (f Filter) Predicate() th.Predicate {
	if f == nil {
		return th.Any()
	}
	return func(ctx context.Context, u telego.Update) bool {
		return f(new(u))
	}
}

// And requires all filters to match.
func And(filters ...Filter) Filter {
	return func(u *telego.Update) bool {
		for _, f := range filters {
			if f != nil && !f(u) {
				return false
			}
		}
		return true
	}
}

// Or requires any filter to match.
func Or(filters ...Filter) Filter {
	return func(u *telego.Update) bool {
		for _, f := range filters {
			if f != nil && f(u) {
				return true
			}
		}
		return false
	}
}

// Not inverts a filter.
func Not(f Filter) Filter {
	return func(u *telego.Update) bool {
		if f == nil {
			return true
		}
		return !f(u)
	}
}

func toPredicates(filters ...Filter) []th.Predicate {
	if len(filters) == 0 {
		return nil
	}
	preds := make([]th.Predicate, 0, len(filters))
	for _, f := range filters {
		preds = append(preds, f.Predicate())
	}
	return preds
}

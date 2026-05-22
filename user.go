package tgx

import (
	"github.com/mymmrac/telego"

	"pyrorhythm.dev/fn/opt"
)

// UserID returns the acting user for an update, if any.
func UserID(u *telego.Update) opt.Of[int64] {
	if u == nil {
		return opt.Nil[int64]()
	}
	switch {
	case u.Message != nil && u.Message.From != nil:
		return opt.Some(u.Message.From.ID)
	case u.CallbackQuery != nil:
		return opt.Some(u.CallbackQuery.From.ID)
	case u.InlineQuery != nil:
		return opt.Some(u.InlineQuery.From.ID)
	case u.ChosenInlineResult != nil:
		return opt.Some(u.ChosenInlineResult.From.ID)
	default:
		return opt.Nil[int64]()
	}
}

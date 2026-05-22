package tgx

import (
	"github.com/mymmrac/telego"

	"pyrorhythm.dev/fn/opt"
)

// NarrowMessage extracts a message update.
func NarrowMessage(u *telego.Update) opt.Of[telego.Message] {
	if u == nil || u.Message == nil {
		return opt.Nil[telego.Message]()
	}
	return opt.SomeAny(*u.Message)
}

// NarrowEditedMessage extracts an edited message update.
func NarrowEditedMessage(u *telego.Update) opt.Of[telego.Message] {
	if u == nil || u.EditedMessage == nil {
		return opt.Nil[telego.Message]()
	}
	return opt.SomeAny(*u.EditedMessage)
}

// NarrowCallbackQuery extracts a callback query update.
func NarrowCallbackQuery(u *telego.Update) opt.Of[telego.CallbackQuery] {
	if u == nil || u.CallbackQuery == nil {
		return opt.Nil[telego.CallbackQuery]()
	}
	return opt.SomeAny(*u.CallbackQuery)
}

// NarrowInlineQuery extracts an inline query update.
func NarrowInlineQuery(u *telego.Update) opt.Of[telego.InlineQuery] {
	if u == nil || u.InlineQuery == nil {
		return opt.Nil[telego.InlineQuery]()
	}
	return opt.SomeAny(*u.InlineQuery)
}

// NarrowChosenInlineResult extracts a chosen inline result update.
func NarrowChosenInlineResult(u *telego.Update) opt.Of[telego.ChosenInlineResult] {
	if u == nil || u.ChosenInlineResult == nil {
		return opt.Nil[telego.ChosenInlineResult]()
	}
	return opt.SomeAny(*u.ChosenInlineResult)
}

package tgx

import "github.com/mymmrac/telego"

func narrowMessage(u *telego.Update) (telego.Message, bool) {
	if u == nil || u.Message == nil {
		return telego.Message{}, false
	}
	return *u.Message, true
}

func narrowEditedMessage(u *telego.Update) (telego.Message, bool) {
	if u == nil || u.EditedMessage == nil {
		return telego.Message{}, false
	}
	return *u.EditedMessage, true
}

func narrowCallbackQuery(u *telego.Update) (telego.CallbackQuery, bool) {
	if u == nil || u.CallbackQuery == nil {
		return telego.CallbackQuery{}, false
	}
	return *u.CallbackQuery, true
}

func narrowInlineQuery(u *telego.Update) (telego.InlineQuery, bool) {
	if u == nil || u.InlineQuery == nil {
		return telego.InlineQuery{}, false
	}
	return *u.InlineQuery, true
}

func narrowChosenInlineResult(u *telego.Update) (telego.ChosenInlineResult, bool) {
	if u == nil || u.ChosenInlineResult == nil {
		return telego.ChosenInlineResult{}, false
	}
	return *u.ChosenInlineResult, true
}

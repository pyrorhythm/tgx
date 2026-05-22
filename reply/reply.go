// Package reply provides res.Of-style reply helpers.
package reply

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"

	"pyrorhythm.dev/fn/res"

	"pyrorhythm.dev/tgx"
)

// TextParams returns send-message params for chatID and text.
func TextParams(chatID int64, text string) res.Of[*telego.SendMessageParams] {
	return res.OKAny(tu.Message(tu.ID(chatID), text))
}

// SendText sends text via ctx and wraps the API result.
func SendText(ctx *tgx.Ctx, chatID int64, text string) res.Of[*telego.Message] {
	return ctx.SendMessage(chatID, text)
}

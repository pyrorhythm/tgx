// Package reply provides optional res.Of-style reply helpers (v1 minimal).
package reply

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"

	"github.com/pyrorhythm/fn/res"

	"pyrorhythm.dev/tgx"
)

// TextParams returns send-message params for chatID and text.
func TextParams(chatID int64, text string) res.Of[*telego.SendMessageParams] {
	return res.OKAny(tu.Message(tu.ID(chatID), text))
}

// SendText sends text via ctx and wraps the API result.
func SendText(ctx *tgx.Ctx, chatID int64, text string) res.Of[*telego.Message] {
	msg, err := ctx.Bot.SendMessage(ctx.Context(), tu.Message(tu.ID(chatID), text))
	return res.FromAny(msg, err)
}

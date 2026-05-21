package tgx

import (
	"context"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"

	"github.com/pyrorhythm/tgx/di"
)

// Ctx wraps telegohandler.Context with tgx helpers.
type Ctx struct {
	TH     *th.Context
	Bot    *telego.Bot
	Update telego.Update
}

// Context returns the underlying context.
func (c *Ctx) Context() context.Context {
	if c.TH != nil {
		return c.TH.Context()
	}
	return context.Background()
}

// PutCtx stores a typed value in the context bag.
func PutCtx[T any](c *Ctx, v T) *Ctx {
	if c.TH == nil {
		return c
	}
	c.TH = c.TH.WithContext(di.Put(c.TH.Context(), v))
	return c
}

// Get reads a typed value from the context bag.
func Get[T any](c *Ctx) (v T, ok bool) {
	o := di.Get[T](c.Context())
	if !o.Valid() {
		return v, false
	}
	return o.Val(), true
}

func wrapCtx(thCtx *th.Context, update telego.Update) *Ctx {
	return &Ctx{
		TH:     thCtx,
		Bot:    thCtx.Bot(),
		Update: update,
	}
}

// SendMessage sends a text message to chatID.
func (c *Ctx) SendMessage(chatID int64, text string) error {
	_, err := c.Bot.SendMessage(c.Context(), tu.Message(tu.ID(chatID), text))
	return err
}

// AnswerCallback answers a callback query.
func (c *Ctx) AnswerCallback(queryID string, text string) error {
	return c.Bot.AnswerCallbackQuery(c.Context(), &telego.AnswerCallbackQueryParams{
		CallbackQueryID: queryID,
		Text:            text,
	})
}

// AnswerInline answers an inline query.
func (c *Ctx) AnswerInline(queryID string, results []telego.InlineQueryResult, cacheTime int) error {
	return c.Bot.AnswerInlineQuery(c.Context(), &telego.AnswerInlineQueryParams{
		InlineQueryID: queryID,
		Results:       results,
		CacheTime:     cacheTime,
	})
}

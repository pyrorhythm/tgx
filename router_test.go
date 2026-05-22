package tgx_test

import (
	"testing"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"

	"pyrorhythm.dev/tgx"
	"pyrorhythm.dev/tgx/filters"
)

func TestRouterRegistersHandler(t *testing.T) {
	bot, _ := telego.NewBot("123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11")
	updates := make(chan telego.Update)
	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		t.Fatal(err)
	}
	dres := tgx.NewDispatcher(bot, updates)
	if !dres.OK() {
		t.Fatal(err)
	}
	_ = bh

	d := dres.Val()

	got := false
	d.Router().OnMessage(func(_ *tgx.Ctx, _ telego.Message) error {
		got = true
		return nil
	}, tgx.WithFilters(filters.Command("start")))

	u := telego.Update{
		UpdateID: 1,
		Message: &telego.Message{
			Text: "/start",
			Chat: telego.Chat{ID: 1, Type: telego.ChatTypePrivate},
		},
	}
	if err := d.BotHandler().BaseGroup().HandleUpdate(t.Context(), bot, u); err != nil {
		t.Fatal(err)
	}
	if !got {
		t.Fatal("handler not invoked")
	}
}

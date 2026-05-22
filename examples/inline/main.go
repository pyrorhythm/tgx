// Inline query example. Set TOKEN or TELEGRAM_BOT_TOKEN.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"

	"pyrorhythm.dev/tgx"
	tgxmw "pyrorhythm.dev/tgx/middleware"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	bot := tgx.NewBotFromEnv()
	if bot.Err() != nil {
		fmt.Println(bot.Err())
		os.Exit(1)
	}

	updates := tgx.UpdatesLongPolling(ctx, bot.Val(), nil)
	if updates.Err() != nil {
		fmt.Println(updates.Err())
		os.Exit(1)
	}

	dres := tgx.NewDispatcher(bot.Val(), updates.Val())
	if !dres.OK() {
		fmt.Println(dres.Err())
		os.Exit(1)
	}

	d := dres.Val()

	d.BotHandler().Use(tgxmw.TelegoRecover())
	d.Router().OnInlineQuery(func(c *tgx.Ctx, q telego.InlineQuery) error {
		result := tu.ResultArticle(q.ID+"-1", "Pick me", tu.TextMessage("You chose: "+q.Query))
		return c.AnswerInline(q.ID, []telego.InlineQueryResult{result}, 300)
	})

	fmt.Println("inline example running")
	if err := d.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

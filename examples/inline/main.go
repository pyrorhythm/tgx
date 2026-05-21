// Inline query example. Set TOKEN or TELEGRAM_BOT_TOKEN.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"

	"github.com/pyrorhythm/tgx"
	tgxmw "github.com/pyrorhythm/tgx/middleware"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	bot, err := tgx.NewBotFromEnv()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	updates, err := tgx.UpdatesLongPolling(ctx, bot, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	d, err := tgx.NewDispatcher(bot, updates)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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

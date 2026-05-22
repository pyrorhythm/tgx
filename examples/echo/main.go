// Echo bot example. Set TOKEN or TELEGRAM_BOT_TOKEN.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/mymmrac/telego"

	"pyrorhythm.dev/tgx"
	"pyrorhythm.dev/tgx/filters"
	tgxmw "pyrorhythm.dev/tgx/middleware"
	"pyrorhythm.dev/tgx/reply"
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
	if !updates.OK() {
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
	d.Router().Use(tgxmw.Recover(), tgxmw.Logger())
	d.Router().OnMessage(func(c *tgx.Ctx, msg telego.Message) error {
		return reply.SendText(c, msg.Chat.ID, msg.Text).Err()
	}, tgx.WithFilters(filters.Command("echo")))

	fmt.Println("echo bot running")
	if err := d.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

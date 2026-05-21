// Echo bot example. Set TOKEN or TELEGRAM_BOT_TOKEN.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"

	"pyrorhythm.dev/tgx"
	"pyrorhythm.dev/tgx/filters"
	tgxmw "pyrorhythm.dev/tgx/middleware"
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
	d.Router().Use(tgxmw.Recover(), tgxmw.Logger())
	d.Router().OnMessage(func(c *tgx.Ctx, msg telego.Message) error {
		_, err := c.Bot.SendMessage(c.Context(), tu.Message(tu.ID(msg.Chat.ID), msg.Text))
		return err
	}, tgx.WithFilters(filters.Command("echo")))

	fmt.Println("echo bot running")
	if err := d.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

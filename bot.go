package tgx

import (
	"context"
	"fmt"
	"os"

	"github.com/mymmrac/telego"

	"pyrorhythm.dev/fn/res"
)

const (
	envTelegramBotToken = "TELEGRAM_BOT_TOKEN"
	envToken            = "TOKEN"
)

// NewBotFromEnv creates a bot from TELEGRAM_BOT_TOKEN or TOKEN.
func NewBotFromEnv() res.Of[*telego.Bot] {
	token := os.Getenv(envTelegramBotToken)
	if token == "" {
		token = os.Getenv(envToken)
	}
	if token == "" {
		return res.Errn[*telego.Bot](fmt.Sprintf("tgx: missing %s or %s", envTelegramBotToken, envToken))
	}
	bot, err := telego.NewBot(token)
	return res.FromAny(bot, err)
}

// UpdatesLongPolling returns the long-polling updates channel.
func UpdatesLongPolling(
	ctx context.Context,
	bot *telego.Bot,
	params *telego.GetUpdatesParams,
) res.Of[<-chan telego.Update] {
	ch, err := bot.UpdatesViaLongPolling(ctx, params)
	return res.FromAny(ch, err)
}

package tgx

import (
	"context"
	"fmt"
	"os"

	"github.com/mymmrac/telego"
)

const (
	envTelegramBotToken = "TELEGRAM_BOT_TOKEN"
	envToken            = "TOKEN"
)

// NewBotFromEnv creates a bot from TELEGRAM_BOT_TOKEN or TOKEN.
func NewBotFromEnv() (*telego.Bot, error) {
	token := os.Getenv(envTelegramBotToken)
	if token == "" {
		token = os.Getenv(envToken)
	}
	if token == "" {
		return nil, fmt.Errorf("tgx: missing %s or %s", envTelegramBotToken, envToken)
	}
	return telego.NewBot(token)
}

// UpdatesLongPolling returns the long-polling updates channel.
func UpdatesLongPolling(ctx context.Context, bot *telego.Bot, params *telego.GetUpdatesParams) (<-chan telego.Update, error) {
	return bot.UpdatesViaLongPolling(ctx, params)
}

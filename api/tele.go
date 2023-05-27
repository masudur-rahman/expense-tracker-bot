package api

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/telebot.v3"
)

func TeleBotRoutes() (*telebot.Bot, error) {
	settings := telebot.Settings{
		Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(settings)
	if err != nil {
		return nil, err
	}

	bot.Handle("/hello", func(ctx telebot.Context) error {
		return ctx.Send(fmt.Sprintf("Hello %v!", ctx.Sender().Username))
	})

	return bot, err
}

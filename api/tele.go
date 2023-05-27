package api

import (
	"os"
	"time"

	"github.com/masudur-rahman/expense-tracker-bot/api/handlers"

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

	bot.Handle("/hello", handlers.Hello)

	bot.Handle("/", handlers.Welcome)

	return bot, err
}

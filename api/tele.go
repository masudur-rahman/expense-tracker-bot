package api

import (
	"os"
	"time"

	"github.com/masudur-rahman/expense-tracker-bot/api/handlers"
	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"gopkg.in/telebot.v3"
)

func TeleBotRoutes(svc *all.Services) (*telebot.Bot, error) {
	settings := telebot.Settings{
		Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	//printer := pkg.NewPrinter(pkg.Options{Style: table.StyleLight, EnableStdout: true})

	bot, err := telebot.NewBot(settings)
	if err != nil {
		return nil, err
	}

	bot.Handle("/", handlers.Welcome)
	bot.Handle("/hello", handlers.Hello)
	bot.Handle("/test", handlers.Test)
	bot.Handle("/new", handlers.AddAccount(svc))
	//bot.Handle("/add", handlers.AddNewExpense(printer, svc))
	//bot.Handle("/list", handlers.ListExpenses(printer, svc))

	return bot, err
}

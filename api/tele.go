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

	bot.Use(masudur_rahman())

	bot.Handle("/", handlers.Welcome)

	bot.Handle(telebot.OnCallback, handlers.Callback)
	bot.Handle(telebot.OnText, handlers.TransactionTextCallback)

	bot.Handle("/new", handlers.New)
	bot.Handle("/newtxn", handlers.NewTransaction)

	bot.Handle("/user", handlers.ListUsers)
	bot.Handle("/balance", handlers.ListAccounts)

	bot.Handle("/list", handlers.ListTransactions)
	bot.Handle("/expense", handlers.ListExpenses)

	bot.Handle("/allsummary", handlers.TransactionSummaryCallback)
	bot.Handle("/summary", handlers.TransactionSummary)
	bot.Handle("/report", handlers.TransactionReportCallback)

	bot.Handle("/cat", handlers.TransactionCategoryCallback)

	bot.Handle("/sync", handlers.SyncSQLiteDatabase)

	return bot, nil
}

func masudur_rahman() telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(ctx telebot.Context) error {
			if ctx.Sender().Username != "masudur_rahman" {
				return ctx.Send("Only allowed user is `masudur_rahman`")
			}

			return next(ctx)
		}
	}
}

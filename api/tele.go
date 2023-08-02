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

	bot.Handle("/new", handlers.AddAccount)
	bot.Handle("/balance", handlers.ListAccounts)

	bot.Handle("/list", handlers.ListTransactions)
	bot.Handle("/expense", handlers.ListExpenses)

	bot.Handle("/allsummary", handlers.TransactionSummaryCallback)
	bot.Handle("/summary", handlers.TransactionSummary)

	bot.Handle("/cat", handlers.ListTransactionCategories)
	bot.Handle("/subcat", handlers.ListTransactionSubcategories)

	// New transaction with flags
	bot.Handle("/txn", handlers.AddNewTransactions)

	// New transaction with callback
	bot.Handle("/newtxn", handlers.NewTransaction)

	bot.Handle("/nuser", handlers.NewUser)
	bot.Handle("/user", handlers.ListUsers)

	return bot, err
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

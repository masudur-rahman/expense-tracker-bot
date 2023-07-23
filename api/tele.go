package api

import (
	"os"
	"time"

	"github.com/masudur-rahman/expense-tracker-bot/api/handlers"
	"github.com/masudur-rahman/expense-tracker-bot/pkg"
	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"github.com/jedib0t/go-pretty/v6/table"
	"gopkg.in/telebot.v3"
)

func TeleBotRoutes(svc *all.Services) (*telebot.Bot, error) {
	settings := telebot.Settings{
		Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	printer := pkg.NewPrinter(pkg.Options{Style: table.StyleColoredBright, EnableStdout: true})

	bot, err := telebot.NewBot(settings)
	if err != nil {
		return nil, err
	}

	bot.Handle("/", handlers.Welcome)

	bot.Handle(telebot.OnCallback, handlers.Callback(svc))
	bot.Handle(telebot.OnText, handlers.TransactionTextCallback(svc))

	bot.Handle("/new", handlers.AddAccount(svc))
	bot.Handle("/balance", handlers.ListAccounts(printer, svc))

	// New transaction with flags
	bot.Handle("/txn", handlers.AddNewTransactions(svc))
	bot.Handle("/list", handlers.ListTransactions(printer, svc))
	bot.Handle("/expense", handlers.ListExpenses(printer, svc))

	bot.Handle("/summary-full", handlers.TransactionSummaryCallback(svc))
	bot.Handle("/summary", handlers.TransactionSummary(printer, svc))

	bot.Handle("/cat", handlers.ListTransactionCategories(svc))
	bot.Handle("/subcat", handlers.ListTransactionSubcategories(svc))

	// New transaction with callback
	bot.Handle("/newtxn", handlers.NewTransaction(svc))

	bot.Handle("/nuser", handlers.NewUser(svc))
	bot.Handle("/user", handlers.ListUsers(svc))

	return bot, err
}

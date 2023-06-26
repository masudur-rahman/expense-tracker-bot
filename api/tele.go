package api

import (
	"fmt"
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
	bot.Handle("/hello", handlers.Hello)
	bot.Handle("/test", handlers.Test)
	bot.Handle(telebot.OnCallback, handlers.Callback(svc))

	bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		fmt.Println(ctx.Text(), "<==>", ctx.Message().Text)
		return ctx.Send("Removing keyboard", &telebot.SendOptions{
			ReplyTo:     ctx.Message(),
			ReplyMarkup: &telebot.ReplyMarkup{RemoveKeyboard: true},
		})
		//return ctx.Reply("You chose " + ctx.Text())
	})

	bot.Handle("/new", handlers.AddAccount(svc))
	bot.Handle("/accounts", handlers.ListAccounts(printer, svc))

	bot.Handle("/txn", handlers.AddNewTransactions(svc))
	bot.Handle("/list", handlers.ListTransactions(printer, svc))

	bot.Handle("/cat", handlers.ListTransactionCategories(svc))
	bot.Handle("/subcat", handlers.ListTransactionSubcategories(svc))

	bot.Handle("/newtxn", handlers.NewTransaction(svc))

	//bot.Handle("/add", handlers.AddNewExpense(printer, svc))
	//bot.Handle("/list", handlers.ListExpenses(printer, svc))

	return bot, err
}

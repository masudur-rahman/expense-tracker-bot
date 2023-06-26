package handlers

import (
	"log"

	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/pkg"
	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"gopkg.in/telebot.v3"
)

func ListTransactionCategories(svc *all.Services) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		txns, err := svc.Txn.ListTxnCategories()
		if err != nil {
			log.Println(err)
			return ctx.Send("Can't list the transaction categories. Please contact the administrator")
		}

		return ctx.Send("Choose one: ", &telebot.SendOptions{
			ReplyTo: ctx.Message(),
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: func() [][]telebot.InlineButton {
					var key []telebot.InlineButton
					for _, cat := range txns {
						key = append(key, telebot.InlineButton{Text: cat.Name, Data: cat.ID})
					}
					return [][]telebot.InlineButton{key}
				}(),
			},
		})
	}
}

func ListTransactionSubcategories(svc *all.Services) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		cat := pkg.SplitString(ctx.Text(), ' ')
		if len(cat) != 2 {
			return ctx.Send("Syntax error")
		}

		subcats, err := svc.Txn.ListTxnSubcategories(cat[1])
		if err != nil {
			log.Println(err)
			return ctx.Send("Can't list the transaction categories. Please contact the administrator")
		}

		return ctx.Send("Choose one: ", &telebot.SendOptions{
			ReplyTo: ctx.Message(),
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: func() [][]telebot.InlineButton {

					var key []telebot.InlineButton
					for _, subcat := range subcats {
						key = append(key, telebot.InlineButton{Text: subcat.Name, Data: subcat.ID})
					}
					return [][]telebot.InlineButton{key}
				}(),
			},
		})
	}
}

func NewTransaction(svc *all.Services) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		return ctx.Send("Type of Transaction: ", &telebot.SendOptions{
			ReplyTo: ctx.Message(),
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{
						telebot.InlineButton{Text: string(models.ExpenseTransaction), Data: string(models.ExpenseTransaction)},
						telebot.InlineButton{Text: string(models.IncomeTransaction), Data: string(models.IncomeTransaction)},
						telebot.InlineButton{Text: string(models.TransferTransaction), Data: string(models.TransferTransaction)},
					},
				},
			},
		})
	}
}

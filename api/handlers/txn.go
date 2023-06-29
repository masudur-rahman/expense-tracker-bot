package handlers

import (
	"fmt"
	"log"

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

		fmt.Println("Subcategory length for", cat[1], ":", len(subcats))

		return ctx.Send("Choose one: ", &telebot.SendOptions{
			ReplyTo: ctx.Message(),
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: func() [][]telebot.InlineButton {
					var keyboard [][]telebot.InlineButton
					var inlineBtn []telebot.InlineButton
					for _, subcat := range subcats {
						inlineBtn = append(inlineBtn, telebot.InlineButton{Text: subcat.Name, Data: subcat.ID})
						if len(inlineBtn) == 3 {
							keyboard = append(keyboard, inlineBtn)
							inlineBtn = nil
						}
					}
					if len(inlineBtn) > 0 {
						keyboard = append(keyboard, inlineBtn)
					}
					return keyboard
				}(),
				ResizeKeyboard: true,
			},
		})
	}
}

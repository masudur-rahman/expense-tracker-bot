package handlers

import (
	"fmt"

	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"gopkg.in/telebot.v3"
)

type CallbackType string

const (
	TxnCategoryType    CallbackType = "txn-category"
	TxnSubcategoryType CallbackType = "txn-subcategory"
)

type CallbackOptions struct {
	Type CallbackType
}

var messageData = make(map[int]string)

func Callback(svc *all.Services) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		fmt.Println(ctx.Callback().Data, "<==>", ctx.Message().Text)
		messageData[ctx.Message().ID] = "Okay, now I've got it..."

		return ctx.Send("Hello there!", &telebot.SendOptions{
			ReplyTo: ctx.Message(),
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{
						telebot.InlineButton{
							Text: "Option 1",
							Data: "Option 1",
						},
					},
				},
				ForceReply:     true,
				ResizeKeyboard: true,
				Placeholder:    "What, Why, How to do it ?",
			},
		})
		//return ctx.EditOrSend("Hello there..! You selected " + ctx.Callback().Data)
	}
}

func TextCallback(svc *all.Services) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {

		fmt.Println(ctx.Text(), "<==>", ctx.Message().Text)
		//return ctx.Send("Removing keyboard", &telebot.SendOptions{
		//	ReplyTo:     ctx.Message(),
		//	ReplyMarkup: &telebot.ReplyMarkup{RemoveKeyboard: true},
		//})

		var data string
		if ctx.Update().Message.ReplyTo != nil {
			replyToID := ctx.Update().Message.ReplyTo.ID
			data = messageData[replyToID]
		}

		return ctx.Reply("You chose "+ctx.Text()+fmt.Sprintf(" [%v]", data), &telebot.SendOptions{
			ReplyTo: ctx.Message(),
		})
	}
}

func Reply() func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		msg, err := ctx.Bot().Reply(ctx.Message(), "Hello there!", &telebot.SendOptions{
			ReplyTo: ctx.Message(),
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{
						telebot.InlineButton{
							Text: "Option 1",
							Data: "Option 1",
						},
					},
				},
				ForceReply:     true,
				ResizeKeyboard: true,
				Placeholder:    "What, Why, How to do it ?",
			},
			ParseMode: "",
			Entities:  nil,
			Protected: false,
		})
		if err != nil {
			return err
		}

		messageData[msg.ID] = "Okay, now I've got it..."
		return nil
	}
}

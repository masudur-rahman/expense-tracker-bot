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

func Callback(svc *all.Services) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		fmt.Println(ctx.Callback().Data, "<==>", ctx.Message().Text)
		return ctx.EditOrSend("Hello there..! You selected " + ctx.Callback().Data)
	}
}

package handlers

import (
	"fmt"

	"gopkg.in/telebot.v3"
)

func Welcome(ctx telebot.Context) error {
	return ctx.Send("Welcome to Expense Tracker !")
}
func Hello(ctx telebot.Context) error {
	return ctx.Send(fmt.Sprintf("Hello %v!", ctx.Sender().Username))
}

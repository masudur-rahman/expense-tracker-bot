package handlers

import (
	"fmt"
	"log"

	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"gopkg.in/telebot.v3"
)

type UserCallbackOptions struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func handleUserCallback(ctx telebot.Context, callbackOpts CallbackOptions) error {
	msg, err := ctx.Bot().Reply(ctx.Message(), `Reply to this Message with the following data

<id> <name> <email(optional)>
i.e.: john "John Doe" john@doe.com
`, &telebot.SendOptions{
		ReplyTo: ctx.Message(),
	})
	if err != nil {
		return err
	}

	callbackData[msg.ID] = callbackOpts
	return nil
}

func processUserCreation(ctx telebot.Context, uop UserCallbackOptions) error {
	if err := all.GetServices().User.CreateUser(&models.User{
		ID:    uop.ID,
		Name:  uop.Name,
		Email: uop.Email,
	}); err != nil {
		log.Println(err)
		return ctx.Send(err.Error())
	}

	return ctx.Send(fmt.Sprintf("New User [%v] added!", uop.Name))
}

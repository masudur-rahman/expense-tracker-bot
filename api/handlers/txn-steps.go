package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"gopkg.in/telebot.v3"
)

func loanOrBorrowTypeTransaction(callbackOpts CallbackOptions) bool {
	return callbackOpts.Transaction.SubcategoryID == "fin-loan" ||
		callbackOpts.Transaction.SubcategoryID == "fin-borrow"
}

func sendTransactionAmountTypeQuery(ctx telebot.Context, svc *all.Services, callbackOpts CallbackOptions) error {
	callbackOpts.Transaction.NextStep = StepAmount
	amounts := []float64{50, 100, 500}
	inlineButtons := make([]telebot.InlineButton, 0, 3)
	for _, amount := range amounts {
		callbackOpts.Transaction.Amount = amount
		btn, err := generateInlineButton(callbackOpts, fmt.Sprintf("%v", amount))
		if err != nil {
			return ctx.Send("Unexpected server error occurred!")
		}

		inlineButtons = append(inlineButtons, btn)
	}

	msg, err := ctx.Bot().Reply(ctx.Message(), "Select an amount or Reply with an amount to this Message", &telebot.SendOptions{
		ReplyTo: ctx.Message(),
		ReplyMarkup: &telebot.ReplyMarkup{
			InlineKeyboard: generateInlineKeyboard(inlineButtons),
			ForceReply:     true,
		},
	})
	if err != nil {
		return err
	}

	callbackData[msg.ID] = callbackOpts
	return nil

}

func sendTransactionSrcTypeQuery(ctx telebot.Context, svc *all.Services, callbackOpts CallbackOptions) error {
	callbackOpts.Transaction.NextStep = StepSrcID
	inlineButtons, err := generateSrcDstTypeInlineButton(svc, callbackOpts, true)
	if err != nil {
		return ctx.Send("Unexpected server error occurred!")
	}

	return ctx.Send(ctx.Message(), "Select Source Account!", &telebot.SendOptions{
		ReplyTo: ctx.Message(),
		ReplyMarkup: &telebot.ReplyMarkup{
			InlineKeyboard: generateInlineKeyboard(inlineButtons),
			ForceReply:     true,
		},
	})
}

func sendTransactionDstTypeQuery(ctx telebot.Context, svc *all.Services, callbackOpts CallbackOptions) error {
	callbackOpts.Transaction.NextStep = StepDstID
	inlineButtons, err := generateSrcDstTypeInlineButton(svc, callbackOpts, false)
	if err != nil {
		return ctx.Send("Unexpected server error occurred!")
	}

	return ctx.Send(ctx.Message(), "Select Destination Account!", &telebot.SendOptions{
		ReplyTo: ctx.Message(),
		ReplyMarkup: &telebot.ReplyMarkup{
			InlineKeyboard: generateInlineKeyboard(inlineButtons),
			ForceReply:     true,
		},
	})
}

func sendTransactionCategoryQuery(ctx telebot.Context, svc *all.Services, callbackOpts CallbackOptions) error {
	callbackOpts.Transaction.NextStep = StepCategory
	cats, err := svc.Txn.ListTxnCategories()
	if err != nil {
		return err
	}

	inlineButtons := make([]telebot.InlineButton, 0, len(cats))
	for _, cat := range cats {
		callbackOpts.Transaction.CategoryID = cat.ID
		btn, err := generateInlineButton(callbackOpts, cat.Name)
		if err != nil {
			return err
		}
		inlineButtons = append(inlineButtons, btn)
	}

	return ctx.Send(ctx.Message(), "Select Transaction category!", &telebot.SendOptions{
		ReplyTo: ctx.Message(),
		ReplyMarkup: &telebot.ReplyMarkup{
			InlineKeyboard: generateInlineKeyboard(inlineButtons),
			ForceReply:     true,
		},
	})
}

func sendTransactionSubcategoryQuery(ctx telebot.Context, svc *all.Services, callbackOpts CallbackOptions) error {
	callbackOpts.Transaction.NextStep = StepSubcategory
	subcats, err := svc.Txn.ListTxnSubcategories(callbackOpts.Transaction.CategoryID)
	if err != nil {
		return err
	}

	inlineButtons := make([]telebot.InlineButton, 0, len(subcats))
	for _, subcat := range subcats {
		callbackOpts.Transaction.SubcategoryID = subcat.ID
		btn, err := generateInlineButton(callbackOpts, subcat.Name)
		if err != nil {
			return err
		}
		inlineButtons = append(inlineButtons, btn)
	}

	return ctx.Send(ctx.Message(), "Select Transaction subcategory!", &telebot.SendOptions{
		ReplyTo: ctx.Message(),
		ReplyMarkup: &telebot.ReplyMarkup{
			InlineKeyboard: generateInlineKeyboard(inlineButtons),
			ForceReply:     true,
		},
	})
}

func sendTransactionUserQuery(ctx telebot.Context, svc *all.Services, callbackOpts CallbackOptions) error {
	callbackOpts.Transaction.NextStep = StepUser
	users, err := svc.User.ListUsers()
	if err != nil {
		return err
	}

	inlineButtons := make([]telebot.InlineButton, 0, len(users))
	for _, user := range users {
		callbackOpts.Transaction.UserID = user.ID
		btn, err := generateInlineButton(callbackOpts, user.Name)
		if err != nil {
			return err
		}
		inlineButtons = append(inlineButtons, btn)
	}

	return ctx.Send(ctx.Message(), "Select the user associated with the Loan/Borrow!", &telebot.SendOptions{
		ReplyTo: ctx.Message(),
		ReplyMarkup: &telebot.ReplyMarkup{
			InlineKeyboard: generateInlineKeyboard(inlineButtons),
			ForceReply:     true,
		},
	})
}

func sendTransactionRemarksQeury(ctx telebot.Context, svc *all.Services, callbackOpts CallbackOptions) error {
	callbackOpts.Transaction.NextStep = StepRemarks
	btn, err := generateInlineButton(callbackOpts, "Done")
	if err != nil {
		return err
	}
	inlineButtons := []telebot.InlineButton{btn}

	return ctx.Send(ctx.Message(), "Complete the transaction by Pressing Done or Reply with Remarks to this message!", &telebot.SendOptions{
		ReplyTo: ctx.Message(),
		ReplyMarkup: &telebot.ReplyMarkup{
			InlineKeyboard: generateInlineKeyboard(inlineButtons),
			ForceReply:     true,
		},
	})
}

func generateSrcDstTypeInlineButton(svc *all.Services, callbackOpts CallbackOptions, src bool) ([]telebot.InlineButton, error) {
	acs, err := svc.Account.ListAccounts()
	if err != nil {
		return nil, err
	}

	inlineButtons := make([]telebot.InlineButton, 0, len(acs))

	var srcOrDst *string
	if src {
		srcOrDst = &callbackOpts.Transaction.SrcID
	} else {
		srcOrDst = &callbackOpts.Transaction.DstID
	}

	for _, ac := range acs {
		*srcOrDst = ac.ID
		btn, err := generateInlineButton(callbackOpts, ac.Name)
		if err != nil {
			return nil, err
		}

		inlineButtons = append(inlineButtons, btn)
	}

	return inlineButtons, nil
}

func generateInlineButton(doc any, btnText string) (telebot.InlineButton, error) {
	jsonData, err := json.Marshal(doc)
	if err != nil {
		return telebot.InlineButton{}, err
	}
	inlnBtn := telebot.InlineButton{
		Text: btnText,
		Data: string(jsonData),
	}
	return inlnBtn, nil
}

func generateInlineKeyboard(inlineButtons []telebot.InlineButton) [][]telebot.InlineButton {
	var keyboard [][]telebot.InlineButton
	var tmpInlnBtns []telebot.InlineButton
	for _, btn := range inlineButtons {
		tmpInlnBtns = append(tmpInlnBtns, btn)
		if len(tmpInlnBtns) == 3 {
			keyboard = append(keyboard, tmpInlnBtns)
			tmpInlnBtns = nil
		}
	}
	if len(tmpInlnBtns) > 0 {
		keyboard = append(keyboard, tmpInlnBtns)
	}
	return keyboard
}

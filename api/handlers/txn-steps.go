package handlers

import (
	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"gopkg.in/telebot.v3"
)

func loanOrBorrowTypeTransaction(callbackOpts CallbackOptions) bool {
	return callbackOpts.Transaction.SubcategoryID == "fin-loan" ||
		callbackOpts.Transaction.SubcategoryID == "fin-borrow"
}

func sendTransactionAmountTypeQuery(ctx telebot.Context, svc *all.Services, callbackOpts CallbackOptions) error {
	callbackOpts.Transaction.NextStep = StepAmount
	inlineButtons, err := generateAmountTypeInlineButton(callbackOpts)
	if err != nil {
		return ctx.Send("Unexpected server error occurred!")
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
	inlineButtons, err := generateTransactionCategoryTypeInlineButton(svc, callbackOpts)
	if err != nil {
		return ctx.Send("Unexpected server error occurred!")
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
	inlineButtons, err := generateTransactionSubcategoryTypeInlineButton(svc, callbackOpts)
	if err != nil {
		return ctx.Send("Unexpected server error occurred!")
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
	inlineButtons, err := generateTransactionUserTypeInlineButton(svc, callbackOpts)
	if err != nil {
		return ctx.Send("Unexpected server error occurred!")
	}

	return ctx.Send(ctx.Message(), "Select the user associated with the Loan/Borrow!", &telebot.SendOptions{
		ReplyTo: ctx.Message(),
		ReplyMarkup: &telebot.ReplyMarkup{
			InlineKeyboard: generateInlineKeyboard(inlineButtons),
			ForceReply:     true,
		},
	})
}

func sendTransactionRemarksQuery(ctx telebot.Context, svc *all.Services, callbackOpts CallbackOptions) error {
	callbackOpts.Transaction.NextStep = StepRemarks
	inlineButtons, err := generateTransactionRemarksTypeInlineButton(svc, callbackOpts)
	if err != nil {
		return ctx.Send("Unexpected server error occurred!")
	}

	msg, err := ctx.Bot().Reply(ctx.Message(), "Complete the transaction by Pressing Done or Reply with Remarks to this message!", &telebot.SendOptions{
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

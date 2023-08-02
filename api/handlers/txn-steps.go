package handlers

import (
	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"gopkg.in/telebot.v3"
)

func loanOrBorrowTypeTransaction(callbackOpts CallbackOptions) bool {
	return callbackOpts.Transaction.SubcategoryID == models.LoanSubcategoryID ||
		callbackOpts.Transaction.SubcategoryID == models.BorrowSubcategoryID
}

func sendTransactionAmountTypeQuery(ctx telebot.Context, callbackOpts CallbackOptions) error {
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

func sendTransactionSrcTypeQuery(ctx telebot.Context, callbackOpts CallbackOptions) error {
	callbackOpts.Transaction.NextStep = StepSrcID
	inlineButtons, err := generateSrcDstTypeInlineButton(callbackOpts, true)
	if err != nil {
		return ctx.Send("Unexpected server error occurred!")
	}

	return ctx.Send("Select Source Account!", &telebot.SendOptions{
		ReplyTo: ctx.Message(),
		ReplyMarkup: &telebot.ReplyMarkup{
			InlineKeyboard: generateInlineKeyboard(inlineButtons),
			ForceReply:     true,
		},
	})
}

func sendTransactionDstTypeQuery(ctx telebot.Context, callbackOpts CallbackOptions) error {
	callbackOpts.Transaction.NextStep = StepDstID
	inlineButtons, err := generateSrcDstTypeInlineButton(callbackOpts, false)
	if err != nil {
		return ctx.Send("Unexpected server error occurred!")
	}

	return ctx.Send("Select Destination Account!", &telebot.SendOptions{
		ReplyTo: ctx.Message(),
		ReplyMarkup: &telebot.ReplyMarkup{
			InlineKeyboard: generateInlineKeyboard(inlineButtons),
			ForceReply:     true,
		},
	})
}

func sendTransactionCategoryQuery(ctx telebot.Context, callbackOpts CallbackOptions) error {
	callbackOpts.Transaction.NextStep = StepCategory
	inlineButtons, err := generateTransactionCategoryTypeInlineButton(callbackOpts)
	if err != nil {
		return ctx.Send("Unexpected server error occurred!")
	}

	return ctx.Send("Select Transaction category!", &telebot.SendOptions{
		ReplyTo: ctx.Message(),
		ReplyMarkup: &telebot.ReplyMarkup{
			InlineKeyboard: generateInlineKeyboard(inlineButtons),
			ForceReply:     true,
		},
	})
}

func sendTransactionSubcategoryQuery(ctx telebot.Context, callbackOpts CallbackOptions) error {
	callbackOpts.Transaction.NextStep = StepSubcategory
	inlineButtons, err := generateTransactionSubcategoryTypeInlineButton(callbackOpts)
	if err != nil {
		return ctx.Send("Unexpected server error occurred!")
	}

	return ctx.Send("Select Transaction subcategory!", &telebot.SendOptions{
		ReplyTo: ctx.Message(),
		ReplyMarkup: &telebot.ReplyMarkup{
			InlineKeyboard: generateInlineKeyboard(inlineButtons),
			ForceReply:     true,
		},
	})
}

func sendTransactionUserQuery(ctx telebot.Context, callbackOpts CallbackOptions) error {
	callbackOpts.Transaction.NextStep = StepUser
	inlineButtons, err := generateTransactionUserTypeInlineButton(callbackOpts)
	if err != nil {
		return ctx.Send("Unexpected server error occurred!")
	}

	return ctx.Send("Select the user associated with the Loan/Borrow!", &telebot.SendOptions{
		ReplyTo: ctx.Message(),
		ReplyMarkup: &telebot.ReplyMarkup{
			InlineKeyboard: generateInlineKeyboard(inlineButtons),
			ForceReply:     true,
		},
	})
}

func sendTransactionRemarksQuery(ctx telebot.Context, callbackOpts CallbackOptions) error {
	callbackOpts.Transaction.NextStep = StepRemarks
	inlineButtons, err := generateTransactionRemarksTypeInlineButton(callbackOpts)
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

func processTransaction(txn TransactionCallbackOptions) error {
	return all.GetServices().Txn.AddTransaction(models.Transaction{
		Amount:        txn.Amount,
		SubcategoryID: txn.SubcategoryID,
		Type:          txn.Type,
		SrcID:         txn.SrcID,
		DstID:         txn.DstID,
		UserID:        txn.UserID,
		Remarks:       txn.Remarks,
	})
}

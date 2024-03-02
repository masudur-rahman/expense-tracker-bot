package handlers

import (
	"fmt"

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
		return ctx.Send(models.ErrCommonResponse(err))
	}

	msg, err := ctx.Bot().Reply(ctx.Message(),
		fmt.Sprintf("%vSelect an amount or Reply with an amount to this Message", callbackOpts.LastSelectedValue),
		commonSendOptions(ctx, inlineButtons),
	)
	if err != nil {
		return err
	}

	callbackData[msg.ID] = callbackOpts
	return nil
}

func sendTransactionSrcTypeQuery(ctx telebot.Context, callbackOpts CallbackOptions) error {
	callbackOpts.Transaction.NextStep = StepSrcID
	inlineButtons, err := generateSrcDstTypeInlineButton(ctx, callbackOpts, true)
	if err != nil {
		return ctx.Send(models.ErrCommonResponse(err))
	}

	return ctx.Send(fmt.Sprintf("%vSelect Source Account:", callbackOpts.LastSelectedValue), commonSendOptions(ctx, inlineButtons))
}

func sendTransactionDstTypeQuery(ctx telebot.Context, callbackOpts CallbackOptions) error {
	callbackOpts.Transaction.NextStep = StepDstID
	inlineButtons, err := generateSrcDstTypeInlineButton(ctx, callbackOpts, false)
	if err != nil {
		return ctx.Send(models.ErrCommonResponse(err))
	}

	return ctx.Send(fmt.Sprintf("%vSelect Destination Account:", callbackOpts.LastSelectedValue), commonSendOptions(ctx, inlineButtons))
}

func sendTransactionCategoryQuery(ctx telebot.Context, callbackOpts CallbackOptions) error {
	callbackOpts.Transaction.NextStep = StepCategory
	inlineButtons, err := generateTransactionCategoryTypeInlineButton(callbackOpts)
	if err != nil {
		return ctx.Send(models.ErrCommonResponse(err))
	}

	return ctx.Send(fmt.Sprintf("%vSelect Transaction category:", callbackOpts.LastSelectedValue), commonSendOptions(ctx, inlineButtons))
}

func sendTransactionSubcategoryQuery(ctx telebot.Context, callbackOpts CallbackOptions) error {
	callbackOpts.Transaction.NextStep = StepSubcategory
	inlineButtons, err := generateTransactionSubcategoryTypeInlineButton(callbackOpts)
	if err != nil {
		return ctx.Send(models.ErrCommonResponse(err))
	}

	return ctx.Send(fmt.Sprintf("%vSelect Transaction subcategory:", callbackOpts.LastSelectedValue), commonSendOptions(ctx, inlineButtons))
}

func sendTransactionUserQuery(ctx telebot.Context, callbackOpts CallbackOptions) error {
	callbackOpts.Transaction.NextStep = StepUser
	inlineButtons, err := generateTransactionUserTypeInlineButton(ctx, callbackOpts)
	if err != nil {
		return ctx.Send(models.ErrCommonResponse(err))
	}

	return ctx.Send(fmt.Sprintf("%vSelect the user associated with the Loan/Borrow:", callbackOpts.LastSelectedValue), commonSendOptions(ctx, inlineButtons))
}

func sendTransactionRemarksQuery(ctx telebot.Context, callbackOpts CallbackOptions) error {
	callbackOpts.Transaction.NextStep = StepRemarks
	inlineButtons, err := generateTransactionRemarksTypeInlineButton(callbackOpts)
	if err != nil {
		return ctx.Send(models.ErrCommonResponse(err))
	}

	msg, err := ctx.Bot().Reply(ctx.Message(), fmt.Sprintf("%vComplete the transaction by Pressing Done or Reply with Remarks to this message:", callbackOpts.LastSelectedValue),
		commonSendOptions(ctx, inlineButtons))
	if err != nil {
		return err
	}

	callbackData[msg.ID] = callbackOpts
	return nil
}

func processTransaction(ctx telebot.Context, txn TransactionCallbackOptions) error {
	user, err := all.GetServices().User.GetUserByTelegramID(ctx.Sender().ID)
	if err != nil {
		return ctx.Send(models.ErrCommonResponse(err))
	}

	return all.GetServices().Txn.AddTransaction(models.Transaction{
		UserID:             user.ID,
		Amount:             txn.Amount,
		SubcategoryID:      txn.SubcategoryID,
		Type:               txn.Type,
		SrcID:              txn.SrcID,
		DstID:              txn.DstID,
		DebtorCreditorName: txn.DebtorCreditorName,
		Remarks:            txn.Remarks,
	})
}

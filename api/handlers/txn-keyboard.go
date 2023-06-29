package handlers

import (
	"fmt"

	"github.com/masudur-rahman/expense-tracker-bot/pkg"
	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"gopkg.in/telebot.v3"
)

func generateAmountTypeInlineButton(callbackOpts CallbackOptions) ([]telebot.InlineButton, error) {
	amounts := []float64{50, 100, 500}
	inlineButtons := make([]telebot.InlineButton, 0, 3)
	for _, amount := range amounts {
		callbackOpts.Transaction.Amount = amount
		btn, err := generateInlineButton(callbackOpts, fmt.Sprintf("%v", amount))
		if err != nil {
			return nil, err
		}

		inlineButtons = append(inlineButtons, btn)
	}

	return inlineButtons, nil
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

func generateTransactionCategoryTypeInlineButton(svc *all.Services, callbackOpts CallbackOptions) ([]telebot.InlineButton, error) {
	cats, err := svc.Txn.ListTxnCategories()
	if err != nil {
		return nil, err
	}

	inlineButtons := make([]telebot.InlineButton, 0, len(cats))
	for _, cat := range cats {
		callbackOpts.Transaction.CategoryID = cat.ID
		btn, err := generateInlineButton(callbackOpts, cat.Name)
		if err != nil {
			return nil, err
		}
		inlineButtons = append(inlineButtons, btn)
	}

	return inlineButtons, nil
}

func generateTransactionSubcategoryTypeInlineButton(svc *all.Services, callbackOpts CallbackOptions) ([]telebot.InlineButton, error) {
	subcats, err := svc.Txn.ListTxnSubcategories(callbackOpts.Transaction.CategoryID)
	if err != nil {
		return nil, err
	}

	inlineButtons := make([]telebot.InlineButton, 0, len(subcats))
	for _, subcat := range subcats {
		callbackOpts.Transaction.SubcategoryID = subcat.ID
		btn, err := generateInlineButton(callbackOpts, subcat.Name)
		if err != nil {
			return nil, err
		}
		inlineButtons = append(inlineButtons, btn)
	}

	return inlineButtons, nil
}

func generateTransactionUserTypeInlineButton(svc *all.Services, callbackOpts CallbackOptions) ([]telebot.InlineButton, error) {
	users, err := svc.User.ListUsers()
	if err != nil {
		return nil, err
	}

	inlineButtons := make([]telebot.InlineButton, 0, len(users))
	for _, user := range users {
		callbackOpts.Transaction.UserID = user.ID
		btn, err := generateInlineButton(callbackOpts, user.Name)
		if err != nil {
			return nil, err
		}
		inlineButtons = append(inlineButtons, btn)
	}

	return inlineButtons, nil
}

func generateTransactionRemarksTypeInlineButton(svc *all.Services, callbackOpts CallbackOptions) ([]telebot.InlineButton, error) {
	btn, err := generateInlineButton(callbackOpts, "Done")
	if err != nil {
		return nil, err
	}
	inlineButtons := []telebot.InlineButton{btn}

	return inlineButtons, nil
}

func generateInlineButton(obj any, btnText string) (telebot.InlineButton, error) {
	data, err := pkg.EncodeToBase64(obj)
	if err != nil {
		return telebot.InlineButton{}, err
	}
	inlnBtn := telebot.InlineButton{
		Text: btnText,
		Data: data,
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

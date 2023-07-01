package handlers

import (
	"fmt"
	"reflect"
	"time"

	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"github.com/patrickmn/go-cache"
	"github.com/rs/xid"
	"gopkg.in/telebot.v3"
)

var c *cache.Cache

func init() {
	c = cache.New(6*time.Hour, 24*time.Hour)
}

func generateAmountTypeInlineButton(callbackOpts CallbackOptions) ([]telebot.InlineButton, error) {
	amounts := []float64{50, 100, 500}
	inlineButtons := make([]telebot.InlineButton, 0, 3)
	for _, amount := range amounts {
		callbackOpts.Transaction.Amount = amount
		btn := generateInlineButton(callbackOpts, fmt.Sprintf("%v", amount))
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
		btn := generateInlineButton(callbackOpts, ac.Name)
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
		btn := generateInlineButton(callbackOpts, cat.Name)
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
		btn := generateInlineButton(callbackOpts, subcat.Name)
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
		btn := generateInlineButton(callbackOpts, user.Name)
		inlineButtons = append(inlineButtons, btn)
	}

	return inlineButtons, nil
}

func generateTransactionRemarksTypeInlineButton(svc *all.Services, callbackOpts CallbackOptions) ([]telebot.InlineButton, error) {
	btn := generateInlineButton(callbackOpts, "Done")
	inlineButtons := []telebot.InlineButton{btn}

	return inlineButtons, nil
}

func storeInlineButtonDataIntoCache(obj any) string {
	id := xid.New().String()
	c.Set(id, obj, 0)
	return id
}

func fetchInlineButtonDataFromCache(data string, obj any) error {
	dd, ok := c.Get(data)
	if !ok {
		return fmt.Errorf("no data found")
	}

	reflect.ValueOf(obj).Elem().Set(reflect.ValueOf(dd))
	c.Delete(data)
	return nil
}

func generateInlineButton(obj any, btnText string) telebot.InlineButton {
	return telebot.InlineButton{
		Text: btnText,
		Data: storeInlineButtonDataIntoCache(obj),
	}
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

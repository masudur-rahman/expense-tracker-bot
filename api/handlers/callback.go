package handlers

import (
	"fmt"

	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"github.com/masudur-rahman/go-oneliners"

	"gopkg.in/telebot.v3"
)

type CallbackType string

type TransactionStep string

const (
	TxnCategoryType    CallbackType = "txn-category"
	TxnSubcategoryType CallbackType = "txn-subcategory"

	TransactionTypeCallback CallbackType = "transaction"

	StepTxnType     TransactionStep = "txn-type"
	StepAmount      TransactionStep = "txn-amount"
	StepSrcID       TransactionStep = "txn-srcid"
	StepDstID       TransactionStep = "txn-dstid"
	StepCategory    TransactionStep = "txn-cat"
	StepSubcategory TransactionStep = "txn-subcat"
	StepUser        TransactionStep = "txn-user"
	StepRemarks     TransactionStep = "txn-remarks"
	StepDone        TransactionStep = "txn-done"
)

type CallbackOptions struct {
	Type        CallbackType
	Transaction TransactionCallbackOptions
}

type TransactionCallbackOptions struct {
	NextStep TransactionStep `json:"nextStep"`

	Type models.TransactionType `json:"type"`

	Amount        float64 `json:"amount"`
	SrcID         string  `json:"srcID"`
	DstID         string  `json:"dstID"`
	CategoryID    string  `json:"catID"`
	SubcategoryID string  `json:"subcatID"`
	UserID        string  `json:"userID"`
	Remarks       string  `json:"remarks"`
}

var messageData = make(map[int]string)

var callbackData = make(map[int]CallbackOptions) // map[messageID]CallbackOptions

func NewTransaction(svc *all.Services) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		callbackOpts := CallbackOptions{
			Type: TransactionTypeCallback,
			Transaction: TransactionCallbackOptions{
				NextStep: StepTxnType,
			},
		}
		types := []models.TransactionType{models.ExpenseTransaction, models.IncomeTransaction, models.TransferTransaction}
		inlineButtons := make([]telebot.InlineButton, 0, 3)
		for _, typ := range types {
			callbackOpts.Transaction.Type = typ
			btn := generateInlineButton(callbackOpts, string(typ))
			inlineButtons = append(inlineButtons, btn)
		}

		return ctx.Send("Select Type of the Transaction:", &telebot.SendOptions{
			ReplyTo: ctx.Message(),
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: generateInlineKeyboard(inlineButtons),
				ForceReply:     true,
			},
		})
	}
}

func TransactionCallback(svc *all.Services) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		callbackOpts, err := parseCallbackOptions(ctx)
		if err != nil {
			return ctx.Send("Invalid data or data expired!")
		}

		oneliners.PrettyJson(callbackOpts, "Callback Options")

		switch callbackOpts.Type {
		case TransactionTypeCallback:
			// Type -> Amount -> SrcID (and/or) DstID -> Category -> Subcategory -> (UserID) -> Remarks
			txn := callbackOpts.Transaction
			switch callbackOpts.Transaction.NextStep {
			case StepTxnType:
				return sendTransactionAmountTypeQuery(ctx, svc, callbackOpts)
			case StepAmount:
				if txn.Type == models.IncomeTransaction {
					return sendTransactionDstTypeQuery(ctx, svc, callbackOpts)
				} else {
					return sendTransactionSrcTypeQuery(ctx, svc, callbackOpts)
				}
			case StepSrcID:
				if txn.Type == models.TransferTransaction {
					return sendTransactionDstTypeQuery(ctx, svc, callbackOpts)
				} else {
					return sendTransactionCategoryQuery(ctx, svc, callbackOpts)
				}
			case StepDstID:
				return sendTransactionCategoryQuery(ctx, svc, callbackOpts)
			case StepCategory:
				return sendTransactionSubcategoryQuery(ctx, svc, callbackOpts)
			case StepSubcategory:
				if loanOrBorrowTypeTransaction(callbackOpts) {
					return sendTransactionUserQuery(ctx, svc, callbackOpts)
				} else {
					return sendTransactionRemarksQuery(ctx, svc, callbackOpts)
				}
			case StepUser:
				return sendTransactionRemarksQuery(ctx, svc, callbackOpts)
			case StepRemarks:
				err = processTransaction(svc, callbackOpts.Transaction)
				if err != nil {
					return ctx.Send(err.Error())
				}
				return ctx.Send("Transaction added successfully!")
			}
		default:
		}

		return nil
	}
}

func parseCallbackOptions(ctx telebot.Context) (CallbackOptions, error) {
	var callbackOpts CallbackOptions
	err := fetchInlineButtonDataFromCache(ctx.Callback().Data, &callbackOpts)
	return callbackOpts, err
}

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

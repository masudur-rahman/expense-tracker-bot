package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"gopkg.in/telebot.v3"
)

type CallbackType string

type TransactionStep string

const (
	TxnCategoryType    CallbackType = "txn-category"
	TxnSubcategoryType CallbackType = "txn-subcategory"

	TransactionTypeCallback CallbackType = "transaction"

	StepTxnType TransactionStep = "txn-type"
	StepAmount  TransactionStep = "txn-amount"
	StepSrcID   TransactionStep = "txn-srcid"
	StepDstID   TransactionStep = "txn-dstid"
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

func TransactionCallback(svc *all.Services) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		callbackOpts, err := parseCallbackOptions(ctx)
		if err != nil {
			return ctx.Send("Callback data parse error!")
		}

		switch callbackOpts.Type {
		case TransactionTypeCallback:
			// Type -> Amount -> SrcID (and/or) DstID -> Category -> Subcategory -> (UserID) -> Remarks
			txn := callbackOpts.Transaction
			switch callbackOpts.Transaction.NextStep {
			case StepTxnType:
				callbackOpts.Transaction.NextStep = StepAmount

				// FIXME: set after sending the message, set to new msg id

				amounts := []float64{50, 100, 500}
				inlnBtns := make([]telebot.InlineButton, 0, 3)
				for _, amount := range amounts {
					callbackOpts.Transaction.Amount = amount
					btn, err := generateInlineButton(callbackOpts, fmt.Sprintf("%v", amount))
					if err != nil {
						return ctx.Send("Unexpected server error occurred!")
					}

					inlnBtns = append(inlnBtns, btn)
				}

				msg, err := ctx.Bot().Reply(ctx.Message(), "Hello there!", &telebot.SendOptions{
					ReplyTo: ctx.Message(),
					ReplyMarkup: &telebot.ReplyMarkup{
						InlineKeyboard: generateInlineKeyboard(inlnBtns),
						ForceReply:     true,
					},
				})
				if err != nil {
					return err
				}

				callbackData[msg.ID] = callbackOpts
				return nil

			case StepAmount:
				acs, err := svc.Account.ListAccounts()
				if err != nil {
					return ctx.Send("Unexpected server error occurred!")
				}

				inlnBtns := make([]telebot.InlineButton, 0, len(acs))
				if txn.Type == models.ExpenseTransaction || txn.Type == models.TransferTransaction {
					callbackOpts.Transaction.NextStep = StepSrcID

					for _, ac := range acs {
						callbackOpts.Transaction.SrcID = ac.ID
						btn, err := generateInlineButton(callbackOpts, ac.Name)
						if err != nil {
							return ctx.Send("Unexpected server error occurred!")
						}

						inlnBtns = append(inlnBtns, btn)

					}
				} else {
					callbackOpts.Transaction.NextStep = StepDstID
					for _, ac := range acs {
						callbackOpts.Transaction.DstID = ac.ID
						btn, err := generateInlineButton(callbackOpts, ac.Name)
						if err != nil {
							return ctx.Send("Unexpected server error occurred!")
						}

						inlnBtns = append(inlnBtns, btn)
					}
				}

				return ctx.Send(ctx.Message(), "Hello there!", &telebot.SendOptions{
					ReplyTo: ctx.Message(),
					ReplyMarkup: &telebot.ReplyMarkup{
						InlineKeyboard: generateInlineKeyboard(inlnBtns),
						ForceReply:     true,
					},
				})
			case StepSrcID:

			}
		default:
		}

		return nil
	}
}

func parseCallbackOptions(ctx telebot.Context) (CallbackOptions, error) {
	var callbackOpts CallbackOptions
	err := json.Unmarshal([]byte(ctx.Callback().Data), &callbackOpts)
	return callbackOpts, err
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

func generateInlineKeyboard(inlnBtns []telebot.InlineButton) [][]telebot.InlineButton {
	var keyboard [][]telebot.InlineButton
	var tmpInlnBtns []telebot.InlineButton
	for _, btn := range inlnBtns {
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

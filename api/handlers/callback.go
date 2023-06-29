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

				msg, err := ctx.Bot().Reply(ctx.Message(), "Hello there!", &telebot.SendOptions{
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

			case StepAmount:
				var isSrc bool
				if txn.Type == models.IncomeTransaction {
					callbackOpts.Transaction.NextStep = StepDstID
				} else {
					callbackOpts.Transaction.NextStep = StepSrcID
					isSrc = true
				}
				inlineButtons, err := generateSrcDstTypeInlineButton(svc, callbackOpts, isSrc)
				if err != nil {
					return ctx.Send("Unexpected server error occurred!")
				}

				return ctx.Send(ctx.Message(), "Hello there!", &telebot.SendOptions{
					ReplyTo: ctx.Message(),
					ReplyMarkup: &telebot.ReplyMarkup{
						InlineKeyboard: generateInlineKeyboard(inlineButtons),
						ForceReply:     true,
					},
				})
			case StepSrcID:
				if txn.Type == models.TransferTransaction {
					callbackOpts.Transaction.NextStep = StepDstID
					inlineButtons, err := generateSrcDstTypeInlineButton(svc, callbackOpts, false)
					if err != nil {
						return ctx.Send("Unexpected server error occurred!")
					}

					return ctx.Send(ctx.Message(), "Hello there!", &telebot.SendOptions{
						ReplyTo: ctx.Message(),
						ReplyMarkup: &telebot.ReplyMarkup{
							InlineKeyboard: generateInlineKeyboard(inlineButtons),
							ForceReply:     true,
						},
					})
				} else {
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

					return ctx.Send(ctx.Message(), "Hello there!", &telebot.SendOptions{
						ReplyTo: ctx.Message(),
						ReplyMarkup: &telebot.ReplyMarkup{
							InlineKeyboard: generateInlineKeyboard(inlineButtons),
							ForceReply:     true,
						},
					})
				}
			case StepDstID:
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

				return ctx.Send(ctx.Message(), "Hello there!", &telebot.SendOptions{
					ReplyTo: ctx.Message(),
					ReplyMarkup: &telebot.ReplyMarkup{
						InlineKeyboard: generateInlineKeyboard(inlineButtons),
						ForceReply:     true,
					},
				})
			case StepCategory:
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

				return ctx.Send(ctx.Message(), "Hello there!", &telebot.SendOptions{
					ReplyTo: ctx.Message(),
					ReplyMarkup: &telebot.ReplyMarkup{
						InlineKeyboard: generateInlineKeyboard(inlineButtons),
						ForceReply:     true,
					},
				})
			case StepSubcategory:
				if callbackOpts.Transaction.SubcategoryID == "fin-loan" ||
					callbackOpts.Transaction.SubcategoryID == "fin-borrow" {
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

					return ctx.Send(ctx.Message(), "Hello there!", &telebot.SendOptions{
						ReplyTo: ctx.Message(),
						ReplyMarkup: &telebot.ReplyMarkup{
							InlineKeyboard: generateInlineKeyboard(inlineButtons),
							ForceReply:     true,
						},
					})
				} else {
					callbackOpts.Transaction.NextStep = StepRemarks
					btn, err := generateInlineButton(callbackOpts, "Done")
					if err != nil {
						return err
					}
					inlineButtons := []telebot.InlineButton{btn}

					return ctx.Send(ctx.Message(), "Hello there!", &telebot.SendOptions{
						ReplyTo: ctx.Message(),
						ReplyMarkup: &telebot.ReplyMarkup{
							InlineKeyboard: generateInlineKeyboard(inlineButtons),
							ForceReply:     true,
						},
					})
				}
			case StepUser:
				callbackOpts.Transaction.NextStep = StepRemarks
				btn, err := generateInlineButton(callbackOpts, "Done")
				if err != nil {
					return err
				}
				inlineButtons := []telebot.InlineButton{btn}

				return ctx.Send(ctx.Message(), "Hello there!", &telebot.SendOptions{
					ReplyTo: ctx.Message(),
					ReplyMarkup: &telebot.ReplyMarkup{
						InlineKeyboard: generateInlineKeyboard(inlineButtons),
						ForceReply:     true,
					},
				})
			case StepRemarks:
				storeTransaction()
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

func storeTransaction() {

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

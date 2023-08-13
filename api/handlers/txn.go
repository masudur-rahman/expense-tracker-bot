package handlers

import (
	"fmt"
	"log"

	"github.com/masudur-rahman/expense-tracker-bot/pkg"
	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"gopkg.in/telebot.v3"
)

const (
	StepCategoryID    NextStep = "cat-id"
	StepSubcategoryID NextStep = "subcat-id"
)

type TxnCategoryCallbackOptions struct {
	NextStep      NextStep `json:"nextStep"`
	CategoryID    string   `json:"categoryID"`
	SubcategoryID string   `json:"subcategoryID"`
}

func handleTransactionWithFlagsCallback(ctx telebot.Context, callbackOpts CallbackOptions) error {
	msg, err := ctx.Bot().Reply(ctx.Message(), `Reply to this Message with the following data


<amount> -t=<type> -s=<subcat> -f=<src> -d=<dst> -u=<user> -r=<remarks>
i.e.: 6666 -t=Expense -s=food-rest -f=cash -r="Coffee with no one"
`, &telebot.SendOptions{
		ReplyTo: ctx.Message(),
	})
	if err != nil {
		return err
	}

	callbackData[msg.ID] = callbackOpts
	return nil
}

func TransactionCategoryCallback(ctx telebot.Context) error {
	callbackOpts := CallbackOptions{
		Type: TxnCategoryTypeCallback,
		Category: TxnCategoryCallbackOptions{
			NextStep: StepCategoryID,
		},
	}
	return sendJustTxnCategoryQuery(ctx, callbackOpts)
}

func handleTransactionCategoryCallback(ctx telebot.Context, callbackOptions CallbackOptions) error {
	cat := callbackOptions.Category
	switch cat.NextStep {
	case StepCategoryID:
		return sendJustTxnSubcategoryQuery(ctx, callbackOptions)
	case StepSubcategoryID:
		return sendTransactionCategoryInformation(ctx, callbackOptions.Category)
	default:
		return ctx.Send("Invalid Step")
	}
}

func sendJustTxnCategoryQuery(ctx telebot.Context, callbackOpts CallbackOptions) error {
	inlineButtons, err := generateJustTxnCategoryTypeInlineButton(callbackOpts)
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

func sendJustTxnSubcategoryQuery(ctx telebot.Context, callbackOpts CallbackOptions) error {
	callbackOpts.Category.NextStep = StepSubcategoryID
	inlineButtons, err := generateJustTxnSubcategoryTypeInlineButton(callbackOpts)
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

func sendTransactionCategoryInformation(ctx telebot.Context, cop TxnCategoryCallbackOptions) error {
	txn := all.GetServices().Txn
	cat, err := txn.GetTxnCategoryName(cop.CategoryID)
	if err != nil {
		return ctx.Send("Unexpected server error occurred!")
	}

	subcat, err := txn.GetTxnSubcategoryName(cop.SubcategoryID)
	if err != nil {
		return ctx.Send("Unexpected server error occurred!")
	}

	return ctx.Send(fmt.Sprintf(`Transaction Category Information:

Category: %v (%v)
Subcategory: %v (%v)
`, cat, cop.CategoryID, subcat, cop.SubcategoryID))
}

func ListTransactionSubcategories(ctx telebot.Context) error {
	cat := pkg.SplitString(ctx.Text(), ' ')
	if len(cat) != 2 {
		return ctx.Send("Syntax error")
	}

	subcats, err := all.GetServices().Txn.ListTxnSubcategories(cat[1])
	if err != nil {
		log.Println(err)
		return ctx.Send("Can't list the transaction categories. Please contact the administrator")
	}

	fmt.Println("Subcategory length for", cat[1], ":", len(subcats))

	return ctx.Send("Choose one: ", &telebot.SendOptions{
		ReplyTo: ctx.Message(),
		ReplyMarkup: &telebot.ReplyMarkup{
			InlineKeyboard: func() [][]telebot.InlineButton {
				var keyboard [][]telebot.InlineButton
				var inlineBtn []telebot.InlineButton
				for _, subcat := range subcats {
					inlineBtn = append(inlineBtn, telebot.InlineButton{Text: subcat.Name, Data: subcat.ID})
					if len(inlineBtn) == 3 {
						keyboard = append(keyboard, inlineBtn)
						inlineBtn = nil
					}
				}
				if len(inlineBtn) > 0 {
					keyboard = append(keyboard, inlineBtn)
				}
				return keyboard
			}(),
			ResizeKeyboard: true,
		},
	})
}

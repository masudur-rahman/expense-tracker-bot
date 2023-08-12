package handlers

import (
	"fmt"
	"strconv"

	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/modules/cache"

	"github.com/masudur-rahman/go-oneliners"

	"gopkg.in/telebot.v3"
)

type CallbackType string

type NextStep string

const (
	TxnCategoryType    CallbackType = "txn-category"
	TxnSubcategoryType CallbackType = "txn-subcategory"

	TransactionTypeCallback CallbackType = "transaction"
	SummaryTypeCallback     CallbackType = "summary"
	ReportTypeCallback      CallbackType = "report"

	StepTxnType     NextStep = "txn-type"
	StepAmount      NextStep = "txn-amount"
	StepSrcID       NextStep = "txn-srcid"
	StepDstID       NextStep = "txn-dstid"
	StepCategory    NextStep = "txn-cat"
	StepSubcategory NextStep = "txn-subcat"
	StepUser        NextStep = "txn-user"
	StepRemarks     NextStep = "txn-remarks"
	StepDone        NextStep = "txn-done"
)

type CallbackOptions struct {
	Type        CallbackType               `json:"type"`
	Transaction TransactionCallbackOptions `json:"transaction"`
	Summary     SummaryCallbackOptions     `json:"summary"`
	Report      ReportCallbackOptions      `json:"report"`
}

type TransactionCallbackOptions struct {
	NextStep NextStep `json:"nextStep"`

	Type models.TransactionType `json:"type"`

	Amount        float64 `json:"amount"`
	SrcID         string  `json:"srcID"`
	DstID         string  `json:"dstID"`
	CategoryID    string  `json:"catID"`
	SubcategoryID string  `json:"subcatID"`
	UserID        string  `json:"userID"`
	Remarks       string  `json:"remarks"`
}

var callbackData = make(map[int]CallbackOptions) // map[messageID]CallbackOptions

func NewTransaction(ctx telebot.Context) error {
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

func Callback(ctx telebot.Context) error {
	callbackOpts, err := parseCallbackOptions(ctx)
	if err != nil {
		return ctx.Send("Invalid data or data expired!")
	}

	oneliners.PrettyJson(callbackOpts, "Callback Options")

	switch callbackOpts.Type {
	case TransactionTypeCallback:
		return handleTransactionCallback(ctx, callbackOpts)
	case SummaryTypeCallback:
		return handleSummaryCallback(ctx, callbackOpts)
	case ReportTypeCallback:
		return handleReportCallback(ctx, callbackOpts)
	default:
		return ctx.Send("Invalid Callback type")
	}
}

func handleTransactionCallback(ctx telebot.Context, callbackOpts CallbackOptions) error {
	// Type -> Amount -> SrcID (and/or) DstID -> Category -> Subcategory -> (UserID) -> Remarks
	txn := callbackOpts.Transaction
	switch txn.NextStep {
	case StepTxnType:
		return sendTransactionAmountTypeQuery(ctx, callbackOpts)
	case StepAmount:
		if txn.Type == models.IncomeTransaction {
			return sendTransactionDstTypeQuery(ctx, callbackOpts)
		} else {
			return sendTransactionSrcTypeQuery(ctx, callbackOpts)
		}
	case StepSrcID:
		if txn.Type == models.TransferTransaction {
			return sendTransactionDstTypeQuery(ctx, callbackOpts)
		} else {
			return sendTransactionCategoryQuery(ctx, callbackOpts)
		}
	case StepDstID:
		return sendTransactionCategoryQuery(ctx, callbackOpts)
	case StepCategory:
		return sendTransactionSubcategoryQuery(ctx, callbackOpts)
	case StepSubcategory:
		if loanOrBorrowTypeTransaction(callbackOpts) {
			return sendTransactionUserQuery(ctx, callbackOpts)
		} else {
			return sendTransactionRemarksQuery(ctx, callbackOpts)
		}
	case StepUser:
		return sendTransactionRemarksQuery(ctx, callbackOpts)
	case StepRemarks:
		err := processTransaction(callbackOpts.Transaction)
		if err != nil {
			return ctx.Send(err.Error())
		}
		return ctx.Send("Transaction added successfully!")
	default:
		return ctx.Send("Invalid Step")
	}
}

func parseCallbackOptions(ctx telebot.Context) (CallbackOptions, error) {
	var callbackOpts CallbackOptions
	err := cache.FetchData(ctx.Callback().Data, &callbackOpts)
	return callbackOpts, err
}

func TransactionTextCallback(ctx telebot.Context) error {
	fmt.Println(ctx.Text(), "<==>", ctx.Message().Text)
	//return ctx.Send("Removing keyboard", &telebot.SendOptions{
	//	ReplyTo:     ctx.Message(),
	//	ReplyMarkup: &telebot.ReplyMarkup{RemoveKeyboard: true},
	//})

	if ctx.Update().Message.ReplyTo == nil {
		return ctx.Reply("Wrong keyword or data")
	}

	replyToID := ctx.Update().Message.ReplyTo.ID
	callbackOpts := callbackData[replyToID]
	if callbackOpts.Type != TransactionTypeCallback {
		return ctx.Reply("Callback must be of Transaction type")
	}

	var err error
	switch callbackOpts.Transaction.NextStep {
	case StepAmount:
		callbackOpts.Transaction.Amount, err = strconv.ParseFloat(ctx.Text(), 64)
		if err != nil {
			return ctx.Reply("Amount parse error")
		}

		return handleTransactionCallback(ctx, callbackOpts)
	case StepRemarks:
		callbackOpts.Transaction.Remarks = ctx.Text()
		return handleTransactionCallback(ctx, callbackOpts)
	default:
		return ctx.Reply("yet to be implemented")
	}
}

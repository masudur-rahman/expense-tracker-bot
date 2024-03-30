package handlers

import (
	"fmt"
	"strconv"

	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/modules/cache"
	"github.com/masudur-rahman/expense-tracker-bot/modules/transaction"
	"github.com/masudur-rahman/expense-tracker-bot/pkg"
	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"github.com/masudur-rahman/go-oneliners"

	"gopkg.in/telebot.v3"
)

type CallbackType string

type NextStep string

const (
	TxnCategoryTypeCallback    CallbackType = "txn-category"
	TxnSubcategoryTypeCallback CallbackType = "txn-subcategory"

	TransactionTypeCallback     CallbackType = "Transaction"
	TransactionFlagTypeCallback CallbackType = "Transaction with flags"
	SummaryTypeCallback         CallbackType = "Summary"
	ReportTypeCallback          CallbackType = "Report"
	AccountTypeCallback         CallbackType = "Account"
	UserTypeCallback            CallbackType = "DebtorsCreditors"

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
	Type              CallbackType               `json:"type"`
	Transaction       TransactionCallbackOptions `json:"transaction,omitempty"`
	Summary           SummaryCallbackOptions     `json:"summary,omitempty"`
	Report            ReportCallbackOptions      `json:"report,omitempty"`
	Account           AccountCallbackOptions     `json:"account,omitempty"`
	User              UserCallbackOptions        `json:"user,omitempty"`
	Category          TxnCategoryCallbackOptions `json:"category,omitempty"`
	LastSelectedValue string
}

type TransactionCallbackOptions struct {
	NextStep NextStep `json:"nextStep"`

	Type models.TransactionType `json:"type"`

	Amount             float64 `json:"amount,omitempty"`
	SrcID              string  `json:"srcID,omitempty"`
	DstID              string  `json:"dstID,omitempty"`
	CategoryID         string  `json:"catID,omitempty"`
	SubcategoryID      string  `json:"subcatID,omitempty"`
	DebtorCreditorName string  `json:"userID,omitempty"`
	Remarks            string  `json:"remarks,omitempty"`
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
		btn := generateInlineButton(callbackOpts, typ)
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
	case TransactionFlagTypeCallback:
		return handleTransactionWithFlagsCallback(ctx, callbackOpts)
	case TxnCategoryTypeCallback:
		return handleTransactionCategoryCallback(ctx, callbackOpts)
	case SummaryTypeCallback:
		return handleSummaryCallback(ctx, callbackOpts)
	case ReportTypeCallback:
		return handleReportCallback(ctx, callbackOpts)
	case AccountTypeCallback:
		return handleAccountCallback(ctx, callbackOpts)
	case UserTypeCallback:
		return handleUserCallback(ctx, callbackOpts)
	default:
		return ctx.Send("Invalid Callback type")
	}
}

func handleTransactionCallback(ctx telebot.Context, callbackOpts CallbackOptions) error {
	// Type -> Amount -> SrcID (and/or) DstID -> Category -> Subcategory -> (DebtorCreditorName) -> Remarks
	txn := callbackOpts.Transaction
	switch txn.NextStep {
	case StepTxnType:
		callbackOpts.LastSelectedValue = fmt.Sprintf("Selected Type: *%v*\n\n", callbackOpts.Transaction.Type)
		return sendTransactionAmountTypeQuery(ctx, callbackOpts)
	case StepAmount:
		callbackOpts.LastSelectedValue = fmt.Sprintf("Selected Amount: *%v*\n\n", callbackOpts.Transaction.Amount)
		if txn.Type == models.IncomeTransaction {
			return sendTransactionDstTypeQuery(ctx, callbackOpts)
		} else {
			return sendTransactionSrcTypeQuery(ctx, callbackOpts)
		}
	case StepSrcID:
		callbackOpts.LastSelectedValue = fmt.Sprintf("Selected Source: *%v*\n\n", callbackOpts.Transaction.SrcID)
		if txn.Type == models.TransferTransaction {
			return sendTransactionDstTypeQuery(ctx, callbackOpts)
		} else {
			return sendTransactionCategoryQuery(ctx, callbackOpts)
		}
	case StepDstID:
		callbackOpts.LastSelectedValue = fmt.Sprintf("Selected Destination: *%v*\n\n", callbackOpts.Transaction.DstID)
		return sendTransactionCategoryQuery(ctx, callbackOpts)
	case StepCategory:
		callbackOpts.LastSelectedValue = fmt.Sprintf("Selected Category: *%v*\n\n", callbackOpts.Transaction.CategoryID)
		return sendTransactionSubcategoryQuery(ctx, callbackOpts)
	case StepSubcategory:
		callbackOpts.LastSelectedValue = fmt.Sprintf("Selected Subcategory: *%v*\n\n", callbackOpts.Transaction.SubcategoryID)
		if loanOrBorrowTypeTransaction(callbackOpts) {
			return sendTransactionUserQuery(ctx, callbackOpts)
		} else {
			return sendTransactionRemarksQuery(ctx, callbackOpts)
		}
	case StepUser:
		callbackOpts.LastSelectedValue = fmt.Sprintf("Selected User: *%v*\n\n", callbackOpts.Transaction.DebtorCreditorName)
		return sendTransactionRemarksQuery(ctx, callbackOpts)
	case StepRemarks:
		err := processTransaction(ctx, callbackOpts.Transaction)
		if err != nil {
			return ctx.Send(err.Error())
		}
		return ctx.Send("Transaction added successfully!")
	default:
		return ctx.Send("Invalid Step")
	}
}

func setLastSelectedValue(callbackOpts *CallbackOptions) {

}

func parseCallbackOptions(ctx telebot.Context) (CallbackOptions, error) {
	var callbackOpts CallbackOptions
	err := cache.FetchData(ctx.Callback().Data, &callbackOpts)
	return callbackOpts, err
}

func TransactionTextCallback(ctx telebot.Context) error {
	if ctx.Update().Message.ReplyTo == nil {
		if err := handleTransactionFromRegularText(ctx); err != nil {
			return ctx.Send(err.Error())
		}
		return ctx.Send("Transaction added successfully!")
	}

	replyToID := ctx.Update().Message.ReplyTo.ID
	callbackOpts := callbackData[replyToID]

	switch callbackOpts.Type {
	case TransactionTypeCallback:
		return handleTransactionTypeTextCallback(ctx, callbackOpts)
	case TransactionFlagTypeCallback:
		return handleTransactionWithFlagTypeTextCallback(ctx, callbackOpts)
	case AccountTypeCallback:
		return handleAccountTypeTextCallback(ctx, callbackOpts)
	case UserTypeCallback:
		return handleUserTypeTextCallback(ctx, callbackOpts)
	default:
		return ctx.Reply("invalid callback type")
	}
}

func handleTransactionTypeTextCallback(ctx telebot.Context, callbackOpts CallbackOptions) error {
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

func handleTransactionWithFlagTypeTextCallback(ctx telebot.Context, callbackOpts CallbackOptions) error {
	var err error
	callbackOpts.Type = TransactionTypeCallback
	callbackOpts.Transaction, err = parseTransactionFlags(ctx.Text())
	if err != nil {
		return ctx.Reply("Data parse error")
	}
	return handleTransactionCallback(ctx, callbackOpts)
}

func handleAccountTypeTextCallback(ctx telebot.Context, callbackOpts CallbackOptions) error {
	switch callbackOpts.Account.NextStep {
	case StepAccountInfo:
		info := pkg.SplitString(ctx.Text(), ' ')
		if len(info) < 2 {
			return ctx.Reply("must contain <id> <account name>")
		}
		callbackOpts.Account.ID, callbackOpts.Account.Name = info[0], info[1]
		return processAccountCreation(ctx, callbackOpts.Account)
	default:
		return ctx.Reply("yet to be implemented")
	}
}

func handleUserTypeTextCallback(ctx telebot.Context, callbackOpts CallbackOptions) error {
	info := pkg.SplitString(ctx.Text(), ' ')
	if len(info) < 2 {
		return ctx.Reply("must contain <id> <name> <email>")
	}
	callbackOpts.User = UserCallbackOptions{
		NickName: info[0],
		FullName: info[1],
		Email: func() string {
			if len(info) > 2 {
				return info[2]
			}
			return ""
		}(),
	}
	return processUserCreation(ctx, callbackOpts.User)
}

func handleTransactionFromRegularText(ctx telebot.Context) error {
	user, err := all.GetServices().User.GetUserByTelegramID(ctx.Sender().ID)
	if err != nil {
		return err
	}

	txn, err := transaction.ParseTransaction(ctx.Text())
	if err != nil {
		return err
	}
	txn.UserID = user.ID
	return all.GetServices().Txn.AddTransaction(txn)
}

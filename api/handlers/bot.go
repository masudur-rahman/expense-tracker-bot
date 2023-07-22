package handlers

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/pkg"
	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/pflag"
	"gopkg.in/telebot.v3"
)

func Welcome(ctx telebot.Context) error {
	return ctx.Send(fmt.Sprintf(`Hello %v %v!
Welcome to Expense Tracker !
Available options are:
/new <type> <unique-name> <Account Name>
/newuser <id> <name> <email>
`, ctx.Sender().FirstName, ctx.Sender().LastName))
}

func ListUsers(svc *all.Services) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		users, err := svc.User.ListUsers()
		if err != nil {
			return ctx.Send(err.Error())
		}

		return ctx.Send(pkg.FormatDocuments(users, "ID", "Name", "Balance"))
	}
}

func NewUser(svc *all.Services) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		// /newuser <id> <name> <email>
		ui := pkg.SplitString(ctx.Text(), ' ')
		if len(ui) < 3 {
			return ctx.Send(`
Syntax unknown.
Format /newuser <id> <name> <email>
`)
		}
		if err := svc.User.CreateUser(&models.User{
			ID:   ui[1],
			Name: ui[2],
			Email: func() string {
				if len(ui) >= 4 {
					return ui[3]
				}
				return ""
			}(),
		}); err != nil {
			log.Println(err)
			return ctx.Send(err.Error())
		}

		return ctx.Send("New User added!")
	}
}

func AddAccount(svc *all.Services) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		// <type (Cash or Bank)> <unique-short-name> <Account Name>
		aci := pkg.SplitString(ctx.Text(), ' ')
		if len(aci) != 4 {
			return ctx.Send(`
Syntax unknown.
Format /new <type> <unique-name> <Account Name>
`)
		}
		acc := &models.Account{
			ID:   aci[2],
			Type: models.AccountType(aci[1]),
			Name: aci[3],
		}
		if err := svc.Account.CreateAccount(acc); err != nil {
			log.Println(err)
			return ctx.Send(err.Error())
		}

		return ctx.Send("New Account Added !")
	}
}

func ListAccounts(printer pkg.Printer, svc *all.Services) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		accounts, err := svc.Account.ListAccounts()
		if err != nil {
			return err
		}

		return ctx.Send(printAccounts(accounts))
	}
}

func printAccounts(accounts []models.Account) string {
	buf := bytes.Buffer{}
	w := tabwriter.NewWriter(&buf, 0, 0, 5, ' ', 0)
	fmt.Fprintln(w, "ID\tType\tName\tBalance")
	for _, ac := range accounts {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", ac.ID, ac.Type, ac.Name, ac.Balance)
	}
	_ = w.Flush()
	return buf.String()
}

/*
If users will be able to select options from the UI, it's ideal to design the input sequence in a way that guides them through the available options. Here's a suggested sequence that facilitates option selection:

1. Type: Ask the user to select the type of transaction (Expense, Income, Transfer). Present the available options as buttons or a dropdown menu.
2. Subcategory: Based on the selected type, present the relevant subcategories as options for the user to choose from. Display them as buttons or in a dropdown menu.
3. Amount: Once the subcategory is selected, prompt the user to enter the monetary amount of the transaction.
4. SrcID/DstID: Depending on the type of transaction, provide the appropriate options for the source ID (for Expense/Transfer) or destination ID (for Income/Transfer). This could be a dropdown menu or a list of selectable options.
5. User (for Loan/Borrow): If the selected subcategory involves a person (Loan or Borrow), present the relevant users as options for the user to select from. Display them as buttons or in a dropdown menu.
6. Remarks: Provide an optional input field for the user to enter any additional remarks or notes related to the transaction.

By structuring the input sequence in this way, users can easily navigate through the available options and make their selections. It enhances the user experience by presenting a guided interface that reduces the chance of errors or confusion during the input process.
*/

type TransactionOptions struct {
	Type     string
	Amount   float64
	SubCatID string
	SrcID    string
	DstID    string
	UserID   string
	Remarks  string
}

func parseTransactionFlags(txnString string) (TransactionOptions, error) {
	var txnOpts TransactionOptions

	set := pflag.NewFlagSet("transaction", pflag.ContinueOnError)
	set.StringVarP(&txnOpts.Type, "type", "t", string(models.ExpenseTransaction), "Type of the transaction")
	//set.Float64VarP(&txnOpts.Amount, "amount", "a", 0, "Transaction amount")
	set.StringVarP(&txnOpts.SubCatID, "subcat", "s", "misc-misc", "Subcategory for the transaction")
	set.StringVarP(&txnOpts.SrcID, "src", "f", "cash", "Source account for the transaction")
	set.StringVarP(&txnOpts.DstID, "dst", "d", "", "Destination account for the transaction")
	set.StringVarP(&txnOpts.UserID, "user", "u", "", "User associated with the loan/borrow")
	set.StringVarP(&txnOpts.Remarks, "remarks", "r", "", "Remarks for the transaction")

	args := strings.Split(txnString, " ")
	err := set.Parse(args)
	if err != nil {
		return TransactionOptions{}, err
	}

	if len(set.Args()) > 0 {
		_, err = fmt.Sscanf(set.Args()[0], "%f", &txnOpts.Amount)
	}

	return txnOpts, err
}

/*
/txn <amount> -t=<type> -s=<subcat> -f=<src> -d=<dst> -u=<user> -r=<remarks>
*/

func AddNewTransactions(svc *all.Services) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		flags := strings.SplitN(ctx.Text(), " ", 2)
		if len(flags) != 2 {
			return ctx.Send("no argument provided for the transaction")
		}

		txnOpts, err := parseTransactionFlags(flags[1])
		if err != nil {
			return ctx.Send(err.Error())
		}
		params := models.Transaction{
			Amount:        txnOpts.Amount,
			SubcategoryID: txnOpts.SubCatID,
			Type:          models.TransactionType(txnOpts.Type),
			SrcID:         txnOpts.SrcID,
			DstID:         txnOpts.DstID,
			UserID:        txnOpts.UserID,
			Timestamp:     time.Now().Unix(),
			Remarks:       txnOpts.Remarks,
		}
		err = svc.Txn.AddTransaction(params)
		if err != nil {
			return ctx.Send(err.Error())
		}

		return ctx.Send("Transaction added")
	}
}

func ListTransactions(printer pkg.Printer, svc *all.Services) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		now := time.Now()
		txns, err := svc.Txn.ListTransactionsByTime(time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Unix(), now.Unix())
		//txns, err := svc.Txn.ListTransactions()
		//txns, err := svc.Txn.ListTransactionsByType(models.ExpenseTransaction)
		if err != nil {
			return err
		}

		//printer.WithRenderType(pkg.RenderTypeMarkdown)
		printer.WithStyle(table.StyleLight)
		printer.WithExceptColumns([]string{"ID"})
		defer printer.ClearColumns()
		printer.PrintDocuments(txns)

		return ctx.Send(pkg.FormatDocuments(txns, "Timestamp", "Amount", "Type"))
		//return ctx.Send(generateTransactionTelegramResponse(txns))
	}
}

package handlers

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/models/gqtypes"
	"github.com/masudur-rahman/expense-tracker-bot/pkg"
	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"github.com/jedib0t/go-pretty/v6/table"
	"gopkg.in/telebot.v3"
)

func Welcome(ctx telebot.Context) error {
	return ctx.Send(`Welcome to Expense Tracker !
Available options are:
/new <type> <unique-name> <Account Name>
`)
}

func Hello(ctx telebot.Context) error {
	return ctx.Send(fmt.Sprintf("Hello %v!", ctx.Sender().Username))
}

func Test(ctx telebot.Context) error {
	return ctx.Send("Choose one: ", &telebot.SendOptions{
		ReplyTo: ctx.Message(),
		ReplyMarkup: &telebot.ReplyMarkup{
			InlineKeyboard: [][]telebot.InlineButton{
				{
					telebot.InlineButton{Text: "Hello", Data: "Hello"},
					telebot.InlineButton{Text: "Bye", Data: "Bye"},
				},
			},
			//ReplyKeyboard: [][]telebot.ReplyButton{
			//	{
			//		telebot.ReplyButton{Text: "Option 1"},
			//		telebot.ReplyButton{Text: "Option 2"},
			//	},
			//	{
			//		telebot.ReplyButton{Text: "Option 3"},
			//		telebot.ReplyButton{Text: "Cancel"},
			//	},
			//},
			//ForceReply:      false,
			//ResizeKeyboard:  true,
			//OneTimeKeyboard: true,
			//RemoveKeyboard:  true,
			//Selective:       false,
			//Placeholder:     "",
		},
		//DisableWebPagePreview: false,
		//DisableNotification:   false,
		//ParseMode:             "",
		//Entities:              nil,
		//AllowWithoutReply:     false,
		//Protected:             false,
	})
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

		//printer.WithColumns("ID", "Type", "Name", "Balance")
		//defer printer.ClearColumns()
		//return ctx.Send(printer.PrintDocuments(accounts))
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

func AddNewTransactions(svc *all.Services) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		// normal expense type
		// <type> <subcat> <amount> <src> <remarks>
		txnOpts := pkg.SplitString(ctx.Text(), ' ')
		if len(txnOpts) < 4 {
			return ctx.Send("Syntax unknown")
		}

		amount, err := strconv.ParseFloat(txnOpts[3], 64)
		if err != nil {
			return ctx.Send("Parsing amount failed")
		}
		var remarks string
		if len(txnOpts) > 5 {
			remarks = txnOpts[5]
		}
		err = svc.Txn.AddTransaction(models.Transaction{
			Amount:        amount,
			SubcategoryID: txnOpts[2],
			Type:          models.TransactionType(txnOpts[1]),
			SrcID:         txnOpts[4],
			Timestamp:     time.Now().Unix(),
			Remarks:       remarks,
		})
		if err != nil {
			return ctx.Send(err.Error())
		}

		return ctx.Send("Transaction added")
	}
}

func ListTransactions(printer pkg.Printer, svc *all.Services) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		txns, err := svc.Txn.ListTransactionsByType(models.ExpenseTransaction)
		if err != nil {
			return err
		}

		//printer.WithRenderType(pkg.RenderTypeMarkdown)
		printer.WithStyle(table.StyleLight)
		printer.WithExceptColumns([]string{"ID"})
		defer printer.ClearColumns()
		printer.PrintDocuments(txns)

		return ctx.Send(pkg.FormatDocuments(txns, "Timestamp", "Amount", "Type", "Remarks"))
		//return ctx.Send(generateTransactionTelegramResponse(txns))
	}
}

func generateTransactionTelegramResponse(txns []models.Transaction) string {
	buf := bytes.Buffer{}
	w := tabwriter.NewWriter(&buf, 0, 0, 5, ' ', 0)
	fmt.Fprintln(w, "Time\tAmount\tType\tDescription")
	for idx := range txns {
		txn := txns[idx]
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", time.Unix(txn.Timestamp, 0).Format("Jan 02, 15:04"), txn.Amount, txn.Type, txn.Remarks)
	}
	_ = w.Flush()
	return buf.String()
}

func AddNewExpense(printer pkg.Printer, svc *all.Services) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		str := pkg.SplitString(ctx.Text(), ' ')
		var err error
		var amount float64
		if len(str) < 3 {
			return ctx.Send(`
Syntax unknown.
Format: /add <amount> <description>
`)
		} else if amount, err = strconv.ParseFloat(str[1], 64); err != nil {
			return ctx.Send(`
Syntax unknown.
Format: /add <amount> <description>
`)
		}

		params := gqtypes.Expense{
			Amount:      amount,
			Description: strings.Join(str[2:], " "),
		}

		printer.PrintDocument(params)
		//if err = svc.Expense.AddExpense(params); err != nil {
		//	return err
		//}

		return ctx.Send(fmt.Sprintf(`
New Expense entry added.
%s: %v Taka
`, params.Description, amount))
	}
}

//func ListExpenses(printer pkg.Printer, svc *all.Services) func(ctx telebot.Context) error {
//	return func(ctx telebot.Context) error {
//		expenses, err := svc.Expense.ListExpenses()
//		if err != nil {
//			return err
//		}
//
//		printer.WithExceptColumns([]string{"ID"})
//		defer printer.ClearColumns()
//		printer.PrintDocuments(expenses)
//
//		return ctx.Send(generateTelegramResponse(expenses))
//	}
//}

func generateTelegramResponse(expenses []*models.Expense) string {
	buf := bytes.Buffer{}
	w := tabwriter.NewWriter(&buf, 0, 0, 5, ' ', 0)
	fmt.Fprintln(w, "Time\tAmount\tDescription")
	for idx := range expenses {
		fmt.Fprintf(w, "%v\t%v\t%v\n", expenses[idx].Date.Format("Jan 02, 15:04"), expenses[idx].Amount, expenses[idx].Description)
	}
	_ = w.Flush()
	return buf.String()
}

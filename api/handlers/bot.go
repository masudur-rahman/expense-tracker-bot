package handlers

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/models/gqtypes"
	"github.com/masudur-rahman/expense-tracker-bot/pkg"
	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"gopkg.in/telebot.v3"
)

func Welcome(ctx telebot.Context) error {
	return ctx.Send("Welcome to Expense Tracker !")
}

func Hello(ctx telebot.Context) error {
	return ctx.Send(fmt.Sprintf("Hello %v!", ctx.Sender().Username))
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
		if err = svc.Expense.AddExpense(params); err != nil {
			return err
		}

		return ctx.Send(fmt.Sprintf(`
New Expense entry added.
%s: %v Taka
`, params.Description, amount))
	}
}

func ListExpenses(printer pkg.Printer, svc *all.Services) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		expenses, err := svc.Expense.ListExpenses()
		if err != nil {
			return err
		}

		printer.WithExceptColumns([]string{"ID"})
		defer printer.ClearColumns()
		printer.PrintDocuments(expenses)

		return ctx.Send(generateTelegramResponse(expenses))
	}
}

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

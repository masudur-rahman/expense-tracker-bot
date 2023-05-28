package handlers

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

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

		out := printer.PrintDocuments(expenses)
		fmt.Printf("\n%v\n", out)
		buf := bytes.Buffer{}
		w := tabwriter.NewWriter(&buf, 0, 0, 5, ' ', 0)
		fmt.Fprintln(w, "Description\tAmount\tTime")
		for idx := range expenses {
			fmt.Fprintf(w, "%s\t%v\t%v\n", expenses[idx].Description, expenses[idx].Amount, expenses[idx].Date.Format(time.RFC822))
		}
		if err = w.Flush(); err != nil {
			return err
		}
		return ctx.Send(buf.String())
	}
}

package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/models/gqtypes"
	"github.com/masudur-rahman/expense-tracker-bot/pkg"
	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"gopkg.in/telebot.v3"
)

type SummaryGroupBy string

type SummaryDuration string

const (
	StepGroupBy  NextStep = "group-by"
	StepDuration NextStep = "duration"

	GroupByTxnType        SummaryGroupBy = "Transaction Type"
	GroupByTxnCategory    SummaryGroupBy = "Category"
	GroupByTxnSubCategory SummaryGroupBy = "Subcategory"

	DurationOneWeek   SummaryDuration = "1 Week"
	DurationThisMonth SummaryDuration = "This Month"
	DurationOneMonth  SummaryDuration = "One Month"
	DurationHalfYear  SummaryDuration = "Last 6 Months"
	DurationThisYear  SummaryDuration = "This Year"
	DurationOneYear   SummaryDuration = "One Year"
)

type SummaryCallbackOptions struct {
	NextStep NextStep `json:"nextStep"`
	GroupBy  SummaryGroupBy
	Duration SummaryDuration
}

func TransactionSummary(printer pkg.Printer, svc *all.Services) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		txns, err := svc.Txn.ListTransactionsByTime("", pkg.StartOfMonth().Unix(), time.Now().Unix())
		if err != nil {
			return err
		}

		summary := generateSummary(txns)
		return ctx.Send(summary.String())
	}
}

func generateSummary(txns []models.Transaction) gqtypes.Summary {
	var summary gqtypes.Summary
	for _, txn := range txns {
		if txn.Type == models.ExpenseTransaction {
			summary.Expense += txn.Amount
		} else {
			summary.Income += txn.Amount
		}
	}
	return summary
}

func TransactionSummaryCallback(svc *all.Services) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		callbackOpts := CallbackOptions{
			Type: SummaryTypeCallback,
			Summary: SummaryCallbackOptions{
				NextStep: StepGroupBy,
			},
		}

		groupBies := []SummaryGroupBy{GroupByTxnType, GroupByTxnCategory, GroupByTxnSubCategory}
		inlineButtons := make([]telebot.InlineButton, 0, 3)
		for _, groupBy := range groupBies {
			callbackOpts.Summary.GroupBy = groupBy
			btn := generateInlineButton(callbackOpts, string(groupBy))
			inlineButtons = append(inlineButtons, btn)
		}

		return ctx.Send("Group By:", &telebot.SendOptions{
			ReplyTo: ctx.Message(),
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: generateInlineKeyboard(inlineButtons),
				ForceReply:     true,
			},
		})
	}
}

func handleSummaryCallback(ctx telebot.Context, callbackOpts CallbackOptions, svc *all.Services) error {
	summary := callbackOpts.Summary
	switch summary.NextStep {
	case StepGroupBy:
		return sendSummaryDurationQuery(ctx, callbackOpts)
	case StepDuration:
		data, err := processSummary(callbackOpts.Summary, svc)
		if err != nil {
			return err
		}

		return ctx.Send(data)
	default:
		return ctx.Send("Invalid Step")
	}
}

func processSummary(smop SummaryCallbackOptions, svc *all.Services) (string, error) {
	now, startTime := time.Now(), pkg.StartOfMonth()
	switch smop.Duration {
	case DurationOneWeek:
		startTime = now.AddDate(0, 0, -7)
	case DurationThisMonth:
		startTime = pkg.StartOfMonth()
	case DurationOneMonth:
		startTime = now.AddDate(0, -1, 0)
	case DurationHalfYear:
		startTime = now.AddDate(0, -6, 0)
	case DurationThisYear:
		startTime = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	case DurationOneYear:
		startTime = now.AddDate(-1, 0, 0)
	}
	txns, err := svc.Txn.ListTransactionsByTime("", startTime.Unix(), now.Unix())
	if err != nil {
		return "", err
	}

	summary := gqtypes.CustomSummary{
		Type:        map[string]gqtypes.FieldCost{},
		Category:    map[string]gqtypes.FieldCost{},
		Subcategory: map[string]gqtypes.FieldCost{},
		Total:       0,
	}
	for _, txn := range txns {
		if smop.GroupBy == GroupByTxnType {
			fc := summary.Type[string(txn.Type)]
			fc.Amount += txn.Amount
			summary.Type[string(txn.Type)] = fc
		} else if smop.GroupBy == GroupByTxnSubCategory {
			fc := summary.Subcategory[txn.SubcategoryID]
			fc.Amount += txn.Amount
			summary.Subcategory[txn.SubcategoryID] = fc
		} else if smop.GroupBy == GroupByTxnCategory {
			cat := strings.Split(txn.SubcategoryID, "-")[0]
			fc := summary.Category[cat]
			fc.Amount += txn.Amount
			summary.Category[cat] = fc
		}

		summary.Total += txn.Amount
	}

	for k, fc := range summary.Category {
		fc.Name, err = svc.Txn.GetTxnCategoryName(k)
		if err != nil {
			return "", err
		}

		summary.Category[k] = fc
	}

	for k, fc := range summary.Subcategory {
		fc.Name, err = svc.Txn.GetTxnSubcategoryName(k)
		if err != nil {
			return "", err
		}

		summary.Subcategory[k] = fc
	}

	return summary.String(), nil
}

func sendSummaryDurationQuery(ctx telebot.Context, callbackOpts CallbackOptions) error {
	callbackOpts.Summary.NextStep = StepDuration
	inlineButtons, err := generateSummaryDurationInlineButton(callbackOpts)
	if err != nil {
		return ctx.Send("Unexpected server error occurred!")
	}

	return ctx.Send("Select Duration for Summary", &telebot.SendOptions{
		ReplyTo: ctx.Message(),
		ReplyMarkup: &telebot.ReplyMarkup{
			InlineKeyboard: generateInlineKeyboard(inlineButtons),
			ForceReply:     true,
		},
	})
}

func generateSummaryDurationInlineButton(callbackOpts CallbackOptions) ([]telebot.InlineButton, error) {
	durations := []SummaryDuration{DurationOneWeek, DurationThisMonth, DurationOneMonth, DurationHalfYear, DurationThisYear, DurationOneYear}
	inlineButtons := make([]telebot.InlineButton, 0, 3)
	for _, duration := range durations {
		callbackOpts.Summary.Duration = duration
		btn := generateInlineButton(callbackOpts, fmt.Sprintf("%v", duration))
		inlineButtons = append(inlineButtons, btn)
	}

	return inlineButtons, nil
}

package handlers

import (
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
	NextStep NextStep        `json:"nextStep"`
	GroupBy  SummaryGroupBy  `json:"groupBy,omitempty"`
	Duration SummaryDuration `json:"duration,omitempty"`
}

func TransactionSummary(ctx telebot.Context) error {
	txns, err := all.GetServices().Txn.ListTransactionsByTime("", pkg.StartOfMonth().Unix(), time.Now().Unix())
	if err != nil {
		return err
	}

	summary := generateSummary(txns)
	return ctx.Send(summary.String())
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

func TransactionSummaryCallback(ctx telebot.Context) error {
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
		btn := generateInlineButton(callbackOpts, groupBy)
		inlineButtons = append(inlineButtons, btn)
	}

	return ctx.Send("Select Summarization Group", &telebot.SendOptions{
		ReplyTo: ctx.Message(),
		ReplyMarkup: &telebot.ReplyMarkup{
			InlineKeyboard: generateInlineKeyboard(inlineButtons),
			ForceReply:     true,
		},
	})
}

func handleSummaryCallback(ctx telebot.Context, callbackOpts CallbackOptions) error {
	summary := callbackOpts.Summary
	switch summary.NextStep {
	case StepGroupBy:
		return sendSummaryDurationQuery(ctx, callbackOpts)
	case StepDuration:
		data, err := processSummary(summary)
		if err != nil {
			return err
		}

		return ctx.Send(data)
	default:
		return ctx.Send("Invalid Step")
	}
}

func processSummary(smop SummaryCallbackOptions) (string, error) {
	now, startTime := time.Now(), calculateStartTime(smop.Duration)

	svc := all.GetServices()
	txns, err := svc.Txn.ListTransactionsByTime("", startTime.Unix(), now.Unix())
	if err != nil {
		return "", err
	}

	summary := gqtypes.SummaryGroups{
		Type:        map[string]gqtypes.FieldCost{},
		Category:    map[string]gqtypes.FieldCost{},
		Subcategory: map[string]gqtypes.FieldCost{},
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
	inlineButtons := generateSummaryDurationInlineButton(callbackOpts)

	return ctx.Send("Select Duration for Summary", &telebot.SendOptions{
		ReplyTo: ctx.Message(),
		ReplyMarkup: &telebot.ReplyMarkup{
			InlineKeyboard: generateInlineKeyboard(inlineButtons),
			ForceReply:     true,
		},
	})
}

func generateSummaryDurationInlineButton(callbackOpts CallbackOptions) []telebot.InlineButton {
	durations := []SummaryDuration{DurationOneWeek, DurationThisMonth, DurationOneMonth, DurationHalfYear, DurationThisYear, DurationOneYear}
	inlineButtons := make([]telebot.InlineButton, 0, 3)
	for _, duration := range durations {
		callbackOpts.Summary.Duration = duration
		btn := generateInlineButton(callbackOpts, duration)
		inlineButtons = append(inlineButtons, btn)
	}

	return inlineButtons
}

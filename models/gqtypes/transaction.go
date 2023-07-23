package gqtypes

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/masudur-rahman/expense-tracker-bot/models"
)

type Transaction struct {
	Amount      float64
	Subcategory string
	Type        models.TransactionType
	Src         string
	Dst         string
	User        string
	Time        string
	Remarks     string
}

type Summary struct {
	Income  float64
	Expense float64
}

func (s Summary) String() string {
	return fmt.Sprintf(`
Transaction Summary

Income:  %v
Expense: %v
`, s.Income, s.Income)
}

type FieldCost struct {
	Name   string
	Amount float64
}

type CustomSummary struct {
	Type        map[string]FieldCost
	Category    map[string]FieldCost
	Subcategory map[string]FieldCost

	Total float64
}

func (s CustomSummary) String() string {
	buf := bytes.Buffer{}
	w := tabwriter.NewWriter(&buf, 0, 0, 5, ' ', 0)
	fmt.Fprintln(w, "Transaction Summary\n")

	for k, v := range s.Type {
		f := v.Name
		if f == "" {
			f = k
		}
		fmt.Fprintln(w, fmt.Sprintf("%v:\t%v", f, v.Amount))
	}

	for k, v := range s.Category {
		f := v.Name
		if f == "" {
			f = k
		}
		fmt.Fprintln(w, fmt.Sprintf("%v:\t%v", f, v.Amount))
	}

	for k, v := range s.Subcategory {
		f := v.Name
		if f == "" {
			f = k
		}
		fmt.Fprintln(w, fmt.Sprintf("%v:\t%v", f, v.Amount))
	}

	fmt.Fprintln(w, fmt.Sprintf("\nTotal:\t%v", s.Total))
	_ = w.Flush()
	return buf.String()
}

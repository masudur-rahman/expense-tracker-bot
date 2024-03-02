package gqtypes

import (
	"bytes"
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

type Transaction struct {
	Date        time.Time
	Type        string
	Amount      float64
	Source      string
	Destination string
	Person      string
	Category    string
	Subcategory string
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

type SummaryGroups struct {
	Type        map[string]FieldCost
	Category    map[string]FieldCost
	Subcategory map[string]FieldCost
}

type Report struct {
	Name         string
	Transactions []Transaction
	Summary      SummaryGroups
}

func (s SummaryGroups) String() string {
	buf := bytes.Buffer{}
	w := tabwriter.NewWriter(&buf, 0, 0, 5, ' ', 0)
	fmt.Fprintln(w, fmt.Sprintf("Transaction Summary\n"))

	for k, v := range s.Type {
		f := v.Name
		if f == "" {
			f = k
		}
		fmt.Fprintln(w, fmt.Sprintf("%v:\t%.2f", f, v.Amount))
	}

	for k, v := range s.Category {
		f := v.Name
		if f == "" {
			f = k
		}
		fmt.Fprintln(w, fmt.Sprintf("%v:\t%.2f", f, v.Amount))
	}

	for k, v := range s.Subcategory {
		f := v.Name
		if f == "" {
			f = k
		}
		fmt.Fprintln(w, fmt.Sprintf("%v:\t%.2f", f, v.Amount))
	}

	_ = w.Flush()
	return buf.String()
}

func (s SummaryGroups) MarkdownString() string {
	buf := bytes.Buffer{}
	w := tabwriter.NewWriter(&buf, 0, 0, 5, ' ', 0)
	fmt.Fprintln(w, "\n## Transaction Summary")

	if len(s.Type) > 0 {
		writeRowHeader(w, "Type", "Amount")
	}
	for k, v := range s.Type {
		f := v.Name
		if f == "" {
			f = k
		}
		fmt.Fprintln(w, fmt.Sprintf("| %v | %v |", f, v.Amount))
	}

	if len(s.Type) > 0 {
		writeRowHeader(w, "Category", "Amount")
	}
	for k, v := range s.Category {
		f := v.Name
		if f == "" {
			f = k
		}
		fmt.Fprintln(w, fmt.Sprintf("| %v | %v |", f, v.Amount))
	}

	if len(s.Type) > 0 {
		writeRowHeader(w, "Subcategory", "Amount")
	}
	for k, v := range s.Subcategory {
		f := v.Name
		if f == "" {
			f = k
		}
		fmt.Fprintln(w, fmt.Sprintf("| %v | %v |", f, v.Amount))
	}

	_ = w.Flush()
	return buf.String()
}

func writeRowHeader(w io.Writer, a, b any) {
	fmt.Fprintln(w)
	fmt.Fprintln(w, fmt.Sprintf("| %v | %v |", a, b))
	fmt.Fprintln(w, fmt.Sprintf("| --- | --- |"))
}

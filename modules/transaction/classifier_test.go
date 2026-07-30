package transaction

import (
	"errors"
	"testing"

	"github.com/masudur-rahman/khorcha-pati/models"
)

func TestNormalizePhrase(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"  Lunch  ", "lunch"},
		{"had a Lunch", "lunch"},
		{"Salary!!", "salary"},
		{"gave back the money", "gave back money"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := normalizePhrase(tt.in); got != tt.want {
			t.Errorf("normalizePhrase(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestLocalClassify(t *testing.T) {
	tests := []struct {
		name     string
		phrase   string
		wantType models.TransactionType
		locked   bool
		want     string
		wantOK   bool
	}{
		{"salary", "got my salary", models.ExpenseTransaction, false, "fin-sal", true},
		{"withdraw", "atm withdraw", models.ExpenseTransaction, false, "fin-with", true},
		{"multiword beats single", "credit card bill", models.ExpenseTransaction, false, "fin-ccpay", true},
		{"dinner", "dinner", models.ExpenseTransaction, false, "food-rest", true},
		{"lunch", "lunch", models.ExpenseTransaction, false, "food-rest", true},
		{"breakfast", "breakfast", models.ExpenseTransaction, false, "food-rest", true},
		{"groceries", "groceries", models.ExpenseTransaction, false, "food-groc", true},
		{"taxi", "taxi", models.ExpenseTransaction, false, "trans-taxi", true},
		{"hospital", "hospital", models.ExpenseTransaction, false, "health-doc", true},
		{"bare water is beverage", "water", models.ExpenseTransaction, false, "food-bev", true},
		{"water bill is utility", "water bill", models.ExpenseTransaction, false, "house-util", true},
		{"bare mobile is electronics", "mobile", models.ExpenseTransaction, false, "shop-elec", true},
		{"mobile recharge is flexi", "mobile recharge", models.ExpenseTransaction, false, "fin-flexi", true},
		{"no match", "asdf qwer zxcv", models.ExpenseTransaction, false, "", false},
		// Locked type constrains candidates to compatible subcategories.
		{"locked income skips expense-only taxi", "rickshaw", models.IncomeTransaction, true, "misc-asset", true},
		{"locked income sold rickshaw", "rickshaw", models.IncomeTransaction, true, "misc-asset", true},
		{"locked income sold gold", "gold", models.IncomeTransaction, true, "fin-gold", true},
		{"locked expense keeps taxi", "taxi", models.ExpenseTransaction, true, "trans-taxi", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := localClassify(tt.phrase, tt.wantType, tt.locked)
			if ok != tt.wantOK {
				t.Fatalf("localClassify(%q) ok = %v, want %v (got %q)", tt.phrase, ok, tt.wantOK, got)
			}
			if ok && got != tt.want {
				t.Errorf("localClassify(%q) = %q, want %q", tt.phrase, got, tt.want)
			}
		})
	}
}

func TestIsRateLimitErr(t *testing.T) {
	tests := []struct {
		err  error
		want bool
	}{
		{nil, false},
		{errors.New("googleapi: Error 429: RESOURCE_EXHAUSTED"), true},
		{errors.New("rate limit exceeded"), true},
		{errors.New("too many requests"), true},
		{errors.New("invalid api key"), false},
	}
	for _, tt := range tests {
		if got := isRateLimitErr(tt.err); got != tt.want {
			t.Errorf("isRateLimitErr(%v) = %v, want %v", tt.err, got, tt.want)
		}
	}
}

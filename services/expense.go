package services

import (
	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/models/gqtypes"
)

type ExpenseService interface {
	AddExpense(params gqtypes.Expense) error
	ListExpenses() ([]*models.Expense, error)
}

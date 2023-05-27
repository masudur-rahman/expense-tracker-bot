package repos

import "github.com/masudur-rahman/expense-tracker-bot/models"

type ExpenseRepository interface {
	GetLastExpense() (*models.Expense, error)
	ListAllExpenses() ([]*models.Expense, error)
	AddNewExpense(expense *models.Expense) error
	DeleteExpense(id string) error
	EditExpense(expense *models.Expense) error
}

package services

import "github.com/masudur-rahman/expense-tracker-bot/models"

type AccountsService interface {
	GetAccountByID(accID string) (*models.Account, error)
	ListAccounts() ([]models.Account, error)
	ListAccountsByType(typ models.AccountType) ([]models.Account, error)
	CreateAccount(account *models.Account) error
	UpdateAccountBalance(accID string, amount float64) error
	DeleteAccount(accID string) error
}

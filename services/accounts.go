package services

import "github.com/masudur-rahman/expense-tracker-bot/models"

type AccountsService interface {
	GetAccountByShortName(userID int64, accID string) (*models.Account, error)
	ListAccounts(userID int64) ([]models.Account, error)
	ListAccountsByType(userID int64, typ models.AccountType) ([]models.Account, error)
	CreateAccount(account *models.Account) error
	UpdateAccountBalance(userID int64, accID string, amount float64) error
	DeleteAccount(userID int64, accID string) error
}

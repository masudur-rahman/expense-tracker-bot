package repos

import "github.com/masudur-rahman/expense-tracker-bot/models"

type AccountsRepository interface {
	GetAccountByID(userID int64, accID string) (*models.Account, error)
	ListAccounts(userID int64) ([]models.Account, error)
	ListAccountsByType(userID int64, typ models.AccountType) ([]models.Account, error)
	AddNewAccount(account *models.Account) error
	UpdateAccountBalance(userID int64, accID string, txnAmount float64) error
	DeleteAccount(userID int64, accID string) error
}

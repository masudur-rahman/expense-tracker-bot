package repos

import "github.com/masudur-rahman/expense-tracker-bot/models"

type AccountsRepository interface {
	GetAccountByID(accID string) (*models.Account, error)
	ListAccounts() ([]models.Account, error)
	ListAccountsByType(typ models.AccountType) ([]models.Account, error)
	AddNewAccount(account *models.Account) error
	UpdateAccountBalance(accID string, txnAmount float64) error
	DeleteAccount(accID string) error
}

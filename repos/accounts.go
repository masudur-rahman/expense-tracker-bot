package repos

import (
	"github.com/masudur-rahman/expense-tracker-bot/models"

	"github.com/masudur-rahman/database"
)

type AccountsRepository interface {
	WithUnitOfWork(uow database.UnitOfWork) AccountsRepository
	GetAccountByShortName(userID int64, shortName string) (*models.Account, error)
	ListAccounts(userID int64) ([]models.Account, error)
	ListAccountsByType(userID int64, typ models.AccountType) ([]models.Account, error)
	AddNewAccount(account *models.Account) error
	UpdateAccountBalance(userID int64, shortName string, txnAmount float64) error
	DeleteAccount(userID int64, shortName string) error
}

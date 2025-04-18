package accounts

import (
	"fmt"

	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/repos"
)

type accountService struct {
	accRepo repos.AccountsRepository
}

func NewAccountService(accRepo repos.AccountsRepository) *accountService {
	return &accountService{accRepo: accRepo}
}

func (as *accountService) GetAccountByShortName(userID int64, shortName string) (*models.Account, error) {
	return as.accRepo.GetAccountByShortName(userID, shortName)
}

func (as *accountService) ListAccounts(userID int64) ([]models.Account, error) {
	return as.accRepo.ListAccounts(userID)
}

func (as *accountService) ListAccountsByType(userID int64, typ models.AccountType) ([]models.Account, error) {
	return as.accRepo.ListAccountsByType(userID, typ)
}

func (as *accountService) CreateAccount(account *models.Account) error {
	if account.UserID == 0 {
		return fmt.Errorf("user-id can't be empty")
	}
	return as.accRepo.AddNewAccount(account)
}

func (as *accountService) UpdateAccountBalance(userID int64, shortName string, amount float64) error {
	return as.accRepo.UpdateAccountBalance(userID, shortName, amount)
}

func (as *accountService) DeleteAccount(userID int64, shortName string) error {
	return as.accRepo.DeleteAccount(userID, shortName)
}

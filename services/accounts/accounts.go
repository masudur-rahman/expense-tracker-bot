package accounts

import (
	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/repos"
)

type accountService struct {
	accRepo repos.AccountsRepository
}

func NewAccountService(accRepo repos.AccountsRepository) *accountService {
	return &accountService{accRepo: accRepo}
}

func (as *accountService) GetAccountByID(accID string) (*models.Account, error) {
	return as.accRepo.GetAccountByID(accID)
}

func (as *accountService) ListAccounts() ([]models.Account, error) {
	return as.accRepo.ListAccounts()
}

func (as *accountService) ListAccountsByType(typ models.AccountType) ([]models.Account, error) {
	return as.accRepo.ListAccountsByType(typ)
}

func (as *accountService) CreateAccount(account *models.Account) error {
	return as.accRepo.AddNewAccount(account)
}

func (as *accountService) UpdateAccountBalance(accID string, amount float64) error {
	return as.accRepo.UpdateAccountBalance(accID, amount)
}

func (as *accountService) DeleteAccount(accID string) error {
	return as.accRepo.DeleteAccount(accID)
}

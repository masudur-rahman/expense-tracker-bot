package accounts

import (
	"time"

	"github.com/masudur-rahman/expense-tracker-bot/infra/logr"
	"github.com/masudur-rahman/expense-tracker-bot/models"

	isql "github.com/masudur-rahman/database/sql"
)

type SQLAccountsRepository struct {
	db     isql.Database
	logger logr.Logger
}

func NewSQLAccountsRepository(db isql.Database, logger logr.Logger) *SQLAccountsRepository {
	return &SQLAccountsRepository{
		db:     db.Table("account"),
		logger: logger,
	}
}

func (a *SQLAccountsRepository) GetAccountByID(accID string) (*models.Account, error) {
	a.logger.Infow("get account by account id", "account id", accID)
	var acc models.Account
	found, err := a.db.ID(accID).FindOne(&acc)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, models.ErrAccountNotFound{AccID: accID}
	}
	return &acc, nil
}

func (a *SQLAccountsRepository) ListAccounts() ([]models.Account, error) {
	a.logger.Infow("list accounts")
	accs := make([]models.Account, 0)
	err := a.db.FindMany(&accs)
	return accs, err
}

func (a *SQLAccountsRepository) ListAccountsByType(typ models.AccountType) ([]models.Account, error) {
	a.logger.Infow("list accounts", "type", typ)
	accs := make([]models.Account, 0)
	err := a.db.FindMany(&accs, models.Account{Type: typ})
	return accs, err
}

func (a *SQLAccountsRepository) AddNewAccount(account *models.Account) error {
	a.logger.Infow("add new account", "name", account.Name)
	_, err := a.GetAccountByID(account.ID)
	if err == nil {
		return models.ErrAccountAlreadyExist{AccID: account.ID}
	} else if !models.IsErrNotFound(err) {
		return err
	}

	_, err = a.db.InsertOne(account)
	return err
}

func (a *SQLAccountsRepository) UpdateAccountBalance(accID string, txnAmount float64) error {
	a.logger.Infow("updating account balance", "account", accID)
	acc, err := a.GetAccountByID(accID)
	if err != nil {
		return err
	}
	acc.Balance += txnAmount
	acc.LastTxnAmount = txnAmount
	acc.LastTxnTimestamp = time.Now().Unix()

	return a.db.ID(acc.ID).MustCols("balance").UpdateOne(acc)
}

func (a *SQLAccountsRepository) DeleteAccount(accID string) error {
	a.logger.Infow("deleting account", "account", accID)
	return a.db.ID(accID).DeleteOne()
}

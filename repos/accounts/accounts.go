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

func (a *SQLAccountsRepository) GetAccountByShortName(userID int64, shortName string) (*models.Account, error) {
	a.logger.Infow("get account by account id", "account id", shortName)
	var acc models.Account
	found, err := a.db.FindOne(&acc, models.Account{ShortName: shortName, UserID: userID})
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, models.ErrAccountNotFound{AccID: shortName}
	}
	return &acc, nil
}

func (a *SQLAccountsRepository) ListAccounts(userID int64) ([]models.Account, error) {
	a.logger.Infow("list accounts")
	accs := make([]models.Account, 0)
	err := a.db.FindMany(&accs, models.Account{UserID: userID})
	return accs, err
}

func (a *SQLAccountsRepository) ListAccountsByType(userID int64, typ models.AccountType) ([]models.Account, error) {
	a.logger.Infow("list accounts", "type", typ)
	accs := make([]models.Account, 0)
	err := a.db.FindMany(&accs, models.Account{UserID: userID, Type: typ})
	return accs, err
}

func (a *SQLAccountsRepository) AddNewAccount(account *models.Account) error {
	a.logger.Infow("add new account", "name", account.Name)
	_, err := a.GetAccountByShortName(account.UserID, account.ShortName)
	if err == nil {
		return models.ErrAccountAlreadyExist{ShortName: account.ShortName}
	} else if !models.IsErrNotFound(err) {
		return err
	}

	_, err = a.db.InsertOne(account)
	return err
}

func (a *SQLAccountsRepository) UpdateAccountBalance(userID int64, shortName string, txnAmount float64) error {
	a.logger.Infow("updating account balance", "account", shortName)
	acc, err := a.GetAccountByShortName(userID, shortName)
	if err != nil {
		return err
	}
	acc.Balance += txnAmount
	acc.LastTxnAmount = txnAmount
	acc.LastTxnTimestamp = time.Now().Unix()

	return a.db.ID(acc.ID).MustCols("balance").UpdateOne(acc)
}

func (a *SQLAccountsRepository) DeleteAccount(userID int64, shortName string) error {
	a.logger.Infow("deleting account", "account", shortName)
	return a.db.DeleteOne(models.Account{ShortName: shortName, UserID: userID})
}

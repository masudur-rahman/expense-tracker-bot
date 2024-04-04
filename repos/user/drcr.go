package user

import (
	"fmt"
	"time"

	"github.com/masudur-rahman/expense-tracker-bot/infra/logr"
	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/repos"

	"github.com/masudur-rahman/styx"
	isql "github.com/masudur-rahman/styx/sql"
)

type SQLDebtorCreditorRepository struct {
	db     isql.Engine
	logger logr.Logger
}

func NewSQLDebtorCreditorRepository(db isql.Engine, logger logr.Logger) *SQLDebtorCreditorRepository {
	return &SQLDebtorCreditorRepository{
		db:     db.Table("debtors_creditors"),
		logger: logger,
	}
}

func (u *SQLDebtorCreditorRepository) WithUnitOfWork(uow styx.UnitOfWork) repos.DebtorCreditorRepository {
	return &SQLDebtorCreditorRepository{
		db:     uow.SQL.Table("debtors_creditors"),
		logger: u.logger,
	}
}

func (u *SQLDebtorCreditorRepository) GetDebtorCreditorByID(id int64) (*models.DebtorsCreditors, error) {
	u.logger.Infow("finding user by id", "id", id)
	var user models.DebtorsCreditors
	found, err := u.db.ID(id).FindOne(&user)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, models.ErrDebtorCreditorNotFound{}
	}
	return &user, nil
}

func (u *SQLDebtorCreditorRepository) GetDebtorCreditorByName(userID int64, name string) (*models.DebtorsCreditors, error) {
	u.logger.Infow("finding debtor-creditor by name", "userid", userID, "name", name)
	filter := models.DebtorsCreditors{
		UserID:   userID,
		NickName: name,
	}
	var drcr models.DebtorsCreditors
	found, err := u.db.FindOne(&drcr, filter)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, models.ErrDebtorCreditorNotFound{UserID: userID, NickName: name}
	}
	return &drcr, nil
}

func (u *SQLDebtorCreditorRepository) UpdateDebtorCreditorBalance(id int64, txnAmount float64) error {
	u.logger.Infow("updating user")
	drcr, err := u.GetDebtorCreditorByID(id)
	if err != nil {
		return err
	}

	drcr.Balance += txnAmount
	drcr.LastTxnTimestamp = time.Now().Unix()

	return u.db.ID(drcr.ID).MustCols("balance").UpdateOne(drcr)
}

func (u *SQLDebtorCreditorRepository) AddNewDebtorCreditor(drcr *models.DebtorsCreditors) error {
	if drcr.UserID == 0 {
		return fmt.Errorf("user-id can't be empty")
	}
	_, err := u.db.InsertOne(drcr)
	return err
}

func (u *SQLDebtorCreditorRepository) ListDebtorCreditors(userID int64) ([]models.DebtorsCreditors, error) {
	u.logger.Infow("listing users")
	users := make([]models.DebtorsCreditors, 0)
	err := u.db.FindMany(&users, models.DebtorsCreditors{UserID: userID})
	return users, err
}

func (u *SQLDebtorCreditorRepository) DeleteDebtorCreditor(id int64) error {
	u.logger.Infow("deleting user", "id", id)
	return u.db.DeleteOne(models.DebtorsCreditors{ID: id})
}

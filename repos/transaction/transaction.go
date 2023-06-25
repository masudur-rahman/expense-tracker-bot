package transaction

import (
	"fmt"

	"github.com/masudur-rahman/expense-tracker-bot/infra/logr"
	"github.com/masudur-rahman/expense-tracker-bot/models"

	isql "github.com/masudur-rahman/database/sql"
)

type SQLTransactionRepository struct {
	db     isql.Database
	logger logr.Logger
}

func NewSQLTransactionRepository(db isql.Database, logger logr.Logger) *SQLTransactionRepository {
	return &SQLTransactionRepository{
		db:     db.Table("transaction"),
		logger: logger,
	}
}

func (t *SQLTransactionRepository) AddTransaction(txn models.Transaction) error {
	t.logger.Infow("inserting new transaction")
	_, err := t.db.InsertOne(txn)
	return err
}

func (t *SQLTransactionRepository) ListTransactions(filter models.Transaction) ([]models.Transaction, error) {
	t.logger.Infow("list transactions")
	txns := make([]models.Transaction, 0)
	err := t.db.FindMany(&txns, filter)
	return txns, err
}

func (t *SQLTransactionRepository) ListTransactionsByCategory(catID string) ([]models.Transaction, error) {
	t.logger.Infow("list transactions by category")
	txns := make([]models.Transaction, 0)
	err := t.db.Where(fmt.Sprintf("sub_category_id LIKE %s%%", catID)).FindMany(&txns)
	return txns, err
}

func (t *SQLTransactionRepository) ListTransactionsByTime(startTime, endTime int64) ([]models.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

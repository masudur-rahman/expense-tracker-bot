package transaction

import (
	"errors"
	"fmt"
	"time"

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
	txn.Timestamp = time.Now().Unix()
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

func (t *SQLTransactionRepository) ListTransactionsByTime(txnType models.TransactionType, startTime, endTime int64) ([]models.Transaction, error) {
	t.logger.Infow("list transactions by time")
	txns := make([]models.Transaction, 0)
	err := t.db.Where(fmt.Sprintf("timestamp >= ? AND timestamp <= ?"), startTime, endTime).FindMany(&txns, models.Transaction{Type: txnType})
	return txns, err
}

func (ts *SQLTransactionRepository) GetTxnCategoryName(catID string) (string, error) {
	cat := models.TxnCategory{}
	has, err := ts.db.Table("txn_category").ID(catID).FindOne(&cat)
	if err != nil {
		return "", err
	} else if !has {
		return "", errors.New("not found")
	}

	return cat.Name, nil
}

func (t *SQLTransactionRepository) ListTxnCategories() ([]models.TxnCategory, error) {
	t.logger.Infow("list transaction category")
	cats := make([]models.TxnCategory, 0)
	err := t.db.Table("txn_category").FindMany(&cats)
	return cats, err
}

func (ts *SQLTransactionRepository) GetTxnSubcategoryName(subcatID string) (string, error) {
	subcat := models.TxnSubcategory{}
	has, err := ts.db.Table("txn_subcategory").ID(subcatID).FindOne(&subcat)
	if err != nil {
		return "", err
	} else if !has {
		return "", errors.New("not found")
	}

	return subcat.Name, nil
}

func (t *SQLTransactionRepository) ListTxnSubcategories(catID string) ([]models.TxnSubcategory, error) {
	t.logger.Infow("list transaction category")
	subcats := make([]models.TxnSubcategory, 0)
	err := t.db.Table("txn_subcategory").FindMany(&subcats, models.TxnSubcategory{CatID: catID})
	return subcats, err
}

func (t *SQLTransactionRepository) UpdateTxnCategories() error {
	for _, cat := range models.TxnCategories {
		if has, err := t.db.Table("txn_category").ID(cat.ID).FindOne(&models.TxnCategory{}); err != nil {
			return err
		} else if has {
			continue
		}

		if _, err := t.db.Table("txn_category").InsertOne(cat); err != nil {
			return err
		}
	}

	for _, subcat := range models.TxnSubcategories {
		if has, err := t.db.Table("txn_subcategory").ID(subcat.ID).FindOne(&models.TxnSubcategory{}); err != nil {
			return err
		} else if has {
			continue
		}
		if _, err := t.db.Table("txn_subcategory").InsertOne(subcat); err != nil {
			return err
		}
	}
	return nil
}

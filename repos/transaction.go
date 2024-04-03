package repos

import (
	"github.com/masudur-rahman/expense-tracker-bot/models"

	"github.com/masudur-rahman/database"
)

type TransactionRepository interface {
	WithUnitOfWork(uow database.UnitOfWork) TransactionRepository
	AddTransaction(txn models.Transaction) error
	ListTransactionsByCategory(userID int64, catID string) ([]models.Transaction, error)
	ListTransactions(filter models.Transaction) ([]models.Transaction, error)
	ListTransactionsByTime(userID int64, txnType models.TransactionType, startTime, endTime int64) ([]models.Transaction, error)

	GetTxnCategoryName(catID string) (string, error)
	ListTxnCategories() ([]models.TxnCategory, error)
	GetTxnSubcategoryName(subcatID string) (string, error)
	ListTxnSubcategories(catID string) ([]models.TxnSubcategory, error)
	UpdateTxnCategories() error
}

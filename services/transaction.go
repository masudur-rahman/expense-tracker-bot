package services

import "github.com/masudur-rahman/expense-tracker-bot/models"

type TransactionService interface {
	AddTransaction(txn models.Transaction) error
	ListTransactions() ([]models.Transaction, error)
	ListTransactionsByType(txnType models.TransactionType) ([]models.Transaction, error)
	ListTransactionsByCategory(catID string) ([]models.Transaction, error)
	ListTransactionsBySubcategory(subcatID string) ([]models.Transaction, error)
	ListTransactionsByTime(txnType models.TransactionType, startTime, endTime int64) ([]models.Transaction, error)
	ListTransactionsBySourceID(srcID string) ([]models.Transaction, error)
	ListTransactionsByDestinationID(dstID string) ([]models.Transaction, error)
	ListTransactionsByUser(username string) ([]models.Transaction, error)

	GetTxnCategoryName(catID string) (string, error)
	ListTxnCategories() ([]models.TxnCategory, error)
	GetTxnSubcategoryName(subcatID string) (string, error)
	ListTxnSubcategories(catID string) ([]models.TxnSubcategory, error)
	UpdateTxnCategories() error
}

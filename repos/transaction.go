package repos

import "github.com/masudur-rahman/expense-tracker-bot/models"

type TransactionRepository interface {
	AddTransaction(txn models.Transaction) error
	ListTransactionsByCategory(catID string) ([]models.Transaction, error)
	ListTransactions(filter models.Transaction) ([]models.Transaction, error)
	ListTransactionsByTime(startTime, endTime int64) ([]models.Transaction, error)
}

package transaction

import (
	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/repos"
)

type txnService struct {
	txnRepo repos.TransactionRepository
}

func NewTxnService(txnRepo repos.TransactionRepository) *txnService {
	return &txnService{txnRepo: txnRepo}
}

func (ts *txnService) AddTransaction(txn models.Transaction) error {
	return ts.txnRepo.AddTransaction(txn)
}

func (ts *txnService) ListTransactionsByType(txnType models.TransactionType) ([]models.Transaction, error) {
	filter := models.Transaction{
		Type: txnType,
	}
	return ts.txnRepo.ListTransactions(filter)
}

func (ts *txnService) ListTransactionsByCategory(catID string) ([]models.Transaction, error) {
	return ts.txnRepo.ListTransactionsByCategory(catID)
}

func (ts *txnService) ListTransactionsBySubcategory(subcatID string) ([]models.Transaction, error) {
	filter := models.Transaction{
		SubcategoryID: subcatID,
	}
	return ts.txnRepo.ListTransactions(filter)
}

func (ts *txnService) ListTransactionsByTime(startTime, endTime int64) ([]models.Transaction, error) {
	return ts.txnRepo.ListTransactionsByTime(startTime, endTime)
}

func (ts *txnService) ListTransactionsBySourceID(srcID string) ([]models.Transaction, error) {
	filter := models.Transaction{
		SrcID: srcID,
	}
	return ts.txnRepo.ListTransactions(filter)
}

func (ts *txnService) ListTransactionsByDestinationID(dstID string) ([]models.Transaction, error) {
	filter := models.Transaction{
		DstID: dstID,
	}
	return ts.txnRepo.ListTransactions(filter)
}

func (ts *txnService) ListTransactionsByUser(username string) ([]models.Transaction, error) {
	filter := models.Transaction{
		User: username,
	}
	return ts.txnRepo.ListTransactions(filter)
}

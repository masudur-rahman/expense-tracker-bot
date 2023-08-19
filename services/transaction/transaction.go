package transaction

import (
	"fmt"

	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/repos"
)

type txnService struct {
	acRepo    repos.AccountsRepository
	userRepo  repos.UserRepository
	txnRepo   repos.TransactionRepository
	eventRepo repos.EventRepository
}

func NewTxnService(acRepo repos.AccountsRepository, userRepo repos.UserRepository, txnRepo repos.TransactionRepository, evRepo repos.EventRepository) *txnService {
	return &txnService{
		acRepo:    acRepo,
		userRepo:  userRepo,
		txnRepo:   txnRepo,
		eventRepo: evRepo,
	}
}

func (ts *txnService) AddTransaction(txn models.Transaction) error {
	if txn.SubcategoryID == "" {
		return fmt.Errorf("subcategory is required")
	}

	switch txn.Type {
	case models.ExpenseTransaction:
		switch txn.SubcategoryID {
		case models.LoanSubcategoryID, models.BorrowReturnSubID:
			if err := ts.userRepo.UpdateUserBalance(txn.UserID, txn.Amount); err != nil {
				return err
			}
		case models.BorrowSubcategoryID, models.LoanRecoverySubID:
			return fmt.Errorf("borrow or loan recovery type expense should be under Income type")
		}
		if err := ts.acRepo.UpdateAccountBalance(txn.SrcID, -txn.Amount); err != nil {
			return err
		}
	case models.IncomeTransaction:
		switch txn.SubcategoryID {
		case models.BorrowSubcategoryID, models.LoanRecoverySubID:
			if err := ts.userRepo.UpdateUserBalance(txn.UserID, -txn.Amount); err != nil {
				return err
			}
		case models.LoanSubcategoryID, models.BorrowReturnSubID:
			return fmt.Errorf("loan or borrow return type expense should be under Expense type")
		}
		if err := ts.acRepo.UpdateAccountBalance(txn.DstID, txn.Amount); err != nil {
			return err
		}
	case models.TransferTransaction:
		if err := ts.acRepo.UpdateAccountBalance(txn.SrcID, -txn.Amount); err != nil {
			return err
		}
		if err := ts.acRepo.UpdateAccountBalance(txn.DstID, txn.Amount); err != nil {
			return err
		}
	}
	return ts.txnRepo.AddTransaction(txn)
}

func (ts *txnService) ListTransactionsByType(txnType models.TransactionType) ([]models.Transaction, error) {
	filter := models.Transaction{
		Type: txnType,
	}
	return ts.txnRepo.ListTransactions(filter)
}

func (ts *txnService) ListTransactions() ([]models.Transaction, error) {
	return ts.txnRepo.ListTransactions(models.Transaction{})
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

func (ts *txnService) ListTransactionsByTime(txnType models.TransactionType, startTime, endTime int64) ([]models.Transaction, error) {
	return ts.txnRepo.ListTransactionsByTime(txnType, startTime, endTime)
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
		UserID: username,
	}
	return ts.txnRepo.ListTransactions(filter)
}

func (ts *txnService) GetTxnCategoryName(catID string) (string, error) {
	return ts.txnRepo.GetTxnCategoryName(catID)
}

func (ts *txnService) ListTxnCategories() ([]models.TxnCategory, error) {
	return ts.txnRepo.ListTxnCategories()
}

func (ts *txnService) GetTxnSubcategoryName(subcatID string) (string, error) {
	return ts.txnRepo.GetTxnSubcategoryName(subcatID)
}

func (ts *txnService) ListTxnSubcategories(catID string) ([]models.TxnSubcategory, error) {
	return ts.txnRepo.ListTxnSubcategories(catID)
}

func (ts *txnService) UpdateTxnCategories() error {
	return ts.txnRepo.UpdateTxnCategories()
}

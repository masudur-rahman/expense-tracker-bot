package convert

import (
	"strings"
	"time"

	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/models/gqtypes"
	"github.com/masudur-rahman/expense-tracker-bot/modules/cache"
	"github.com/masudur-rahman/expense-tracker-bot/services/all"
)

func ToTransactionAPIFormat(txn models.Transaction) gqtypes.Transaction {
	svc := all.GetServices()
	var err error
	var category, subcategory, src, dst, person string
	catID := strings.Split(txn.SubcategoryID, "-")[0]
	if err = cache.FetchDataWithCustomFunc(catID, &category, func() (any, error) {
		return svc.Txn.GetTxnCategoryName(catID)
	}); err != nil {
		category = catID
	}

	if err = cache.FetchDataWithCustomFunc(txn.SubcategoryID, &subcategory, func() (any, error) {
		return svc.Txn.GetTxnSubcategoryName(txn.SubcategoryID)
	}); err != nil {
		subcategory = txn.SubcategoryID
	}

	if txn.SrcID != "" {
		if err = cache.FetchDataWithCustomFunc(txn.SrcID, &src, func() (any, error) {
			ac, err := svc.Account.GetAccountByID(txn.SrcID)
			if err != nil {
				return nil, err
			}
			return ac.Name, nil
		}); err != nil {
			src = txn.SrcID
		}
	}

	if txn.DstID != "" {
		if err = cache.FetchDataWithCustomFunc(txn.DstID, &dst, func() (any, error) {
			ac, err := svc.Account.GetAccountByID(txn.DstID)
			if err != nil {
				return nil, err
			}
			return ac.Name, nil
		}); err != nil {
			dst = txn.DstID
		}
	}

	if txn.UserID != "" {
		if err = cache.FetchDataWithCustomFunc(txn.UserID, &person, func() (any, error) {
			user, err := svc.User.GetUserByID(txn.UserID)
			if err != nil {
				return nil, err
			}
			return user.Name, nil
		}); err != nil {
			person = txn.UserID
		}
	}

	return gqtypes.Transaction{
		Date:        time.Unix(txn.Timestamp, 0),
		Type:        string(txn.Type),
		Amount:      txn.Amount,
		Source:      src,
		Destination: dst,
		Person:      person,
		Category:    category,
		Subcategory: subcategory,
		Remarks:     txn.Remarks,
	}
}

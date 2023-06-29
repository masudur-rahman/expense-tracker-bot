package all

import (
	"github.com/masudur-rahman/expense-tracker-bot/infra/logr"
	"github.com/masudur-rahman/expense-tracker-bot/repos/accounts"
	"github.com/masudur-rahman/expense-tracker-bot/repos/event"
	"github.com/masudur-rahman/expense-tracker-bot/repos/transaction"
	"github.com/masudur-rahman/expense-tracker-bot/repos/user"
	"github.com/masudur-rahman/expense-tracker-bot/services"
	accsvc "github.com/masudur-rahman/expense-tracker-bot/services/accounts"
	eventsvc "github.com/masudur-rahman/expense-tracker-bot/services/event"
	txnsvc "github.com/masudur-rahman/expense-tracker-bot/services/transaction"
	usersvc "github.com/masudur-rahman/expense-tracker-bot/services/user"

	isql "github.com/masudur-rahman/database/sql"
)

type Services struct {
	Account services.AccountsService
	User    services.UserService
	Txn     services.TransactionService
	Event   services.EventService
}

func GetSQLServices(db isql.Database, logger logr.Logger) *Services {
	accRepo := accounts.NewSQLAccountsRepository(db, logger)
	userRepo := user.NewSQLUserRepository(db, logger)
	txnRepo := transaction.NewSQLTransactionRepository(db, logger)
	eventRepo := event.NewSQLEventRepository(db, logger)

	accSvc := accsvc.NewAccountService(accRepo)
	userSvc := usersvc.NewUserService(userRepo)
	txnSvc := txnsvc.NewTxnService(accRepo, userRepo, txnRepo, eventRepo)
	eventSvc := eventsvc.NewEventService(eventRepo)

	return &Services{
		Account: accSvc,
		User:    userSvc,
		Txn:     txnSvc,
		Event:   eventSvc,
	}
}

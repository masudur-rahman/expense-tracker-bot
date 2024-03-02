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
	User           services.UserService
	Account        services.AccountsService
	DebtorCreditor services.DebtorCreditorService
	Txn            services.TransactionService
	Event          services.EventService
}

var svc *Services

func GetServices() *Services {
	return svc
}

func InitiateSQLServices(db isql.Database, logger logr.Logger) {
	userRepo := user.NewSQLUserRepository(db, logger)
	accRepo := accounts.NewSQLAccountsRepository(db, logger)
	drCrRepo := user.NewSQLDebtorCreditorRepository(db, logger)
	txnRepo := transaction.NewSQLTransactionRepository(db, logger)
	eventRepo := event.NewSQLEventRepository(db, logger)

	userSvc := usersvc.NewUserService(userRepo)
	accSvc := accsvc.NewAccountService(accRepo)
	drCrSvc := usersvc.NewDebtorCreditorService(drCrRepo)
	txnSvc := txnsvc.NewTxnService(accRepo, drCrRepo, txnRepo, eventRepo)
	eventSvc := eventsvc.NewEventService(eventRepo)

	svc = &Services{
		User:           userSvc,
		Account:        accSvc,
		DebtorCreditor: drCrSvc,
		Txn:            txnSvc,
		Event:          eventSvc,
	}
}

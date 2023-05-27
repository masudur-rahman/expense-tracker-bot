package all

import (
	"github.com/masudur-rahman/expense-tracker-bot/infra/database/nosql"
	isql "github.com/masudur-rahman/expense-tracker-bot/infra/database/sql"
	"github.com/masudur-rahman/expense-tracker-bot/infra/logr"
	"github.com/masudur-rahman/expense-tracker-bot/repos/expense"
	"github.com/masudur-rahman/expense-tracker-bot/repos/user"
	"github.com/masudur-rahman/expense-tracker-bot/services"
	expensesvc "github.com/masudur-rahman/expense-tracker-bot/services/expense"
	usersvc "github.com/masudur-rahman/expense-tracker-bot/services/user"
)

type Services struct {
	User    services.UserService
	Expense services.ExpenseService
}

func GetNoSQLServices(db nosql.Database, logger logr.Logger) *Services {
	userRepo := user.NewNoSQLUserRepository(db, logger)

	userSvc := usersvc.NewUserService(userRepo)

	return &Services{
		User: userSvc,
	}
}

func GetSQLServices(db isql.Database, logger logr.Logger) *Services {
	userRepo := user.NewSQLUserRepository(db, logger)
	expenseRepo := expense.NewSQLExpenseRepository(db, logger)

	userSvc := usersvc.NewUserService(userRepo)
	expenseSvc := expensesvc.NewExpenseService(expenseRepo)

	return &Services{
		User:    userSvc,
		Expense: expenseSvc,
	}
}

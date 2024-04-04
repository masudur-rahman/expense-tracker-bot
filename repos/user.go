package repos

import (
	"github.com/masudur-rahman/expense-tracker-bot/models"

	"github.com/masudur-rahman/database"
)

type DebtorCreditorRepository interface {
	WithUnitOfWork(uow database.UnitOfWork) DebtorCreditorRepository
	GetDebtorCreditorByID(id int64) (*models.DebtorsCreditors, error)
	GetDebtorCreditorByName(userID int64, name string) (*models.DebtorsCreditors, error)
	ListDebtorCreditors(userID int64) ([]models.DebtorsCreditors, error)
	AddNewDebtorCreditor(drcr *models.DebtorsCreditors) error
	UpdateDebtorCreditorBalance(id int64, amount float64) error
	DeleteDebtorCreditor(id int64) error
}

type UserRepository interface {
	GetUserByID(id int64) (*models.User, error)
	GetUser(filter models.User) (*models.User, error)
	ListUsers() ([]models.User, error)
	AddNewUser(user *models.User) error
	UpdateUser(id int64, user *models.User) error
	DeleteUser(id int64) error
}

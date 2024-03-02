package services

import (
	"github.com/masudur-rahman/expense-tracker-bot/models"
)

type DebtorCreditorService interface {
	GetDebtorCreditorByID(id int64) (*models.DebtorsCreditors, error)
	GetDebtorCreditorByName(userID int64, name string) (*models.DebtorsCreditors, error)
	ListDebtorCreditors(userID int64) ([]models.DebtorsCreditors, error)
	CreateDebtorCreditor(drcr *models.DebtorsCreditors) error
	UpdateDebtorCreditorBalance(id int64, amount float64) error
	DeleteDebtorCreditor(id int64) error
}

type UserService interface {
	GetUserByID(id int64) (*models.User, error)
	GetUserByTelegramID(id int64) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	ListUsers() ([]models.User, error)
	SignUp(user *models.User) error
	UpdateUser(id int64, user *models.User) error
	DeleteUser(id int64) error
}

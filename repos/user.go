package repos

import "github.com/masudur-rahman/expense-tracker-bot/models"

type UserRepository interface {
	GetUserByID(id string) (*models.User, error)
	GetUserByName(username string) (*models.User, error)
	ListUsers() ([]models.User, error)
	AddNewUser(user *models.User) error
	UpdateUserBalance(id string, amount float64) error
	DeleteUser(id string) error
}

package services

import (
	"github.com/masudur-rahman/expense-tracker-bot/models"
)

type UserService interface {
	GetUserByID(id string) (*models.User, error)
	GetUserByName(username string) (*models.User, error)
	ListUsers() ([]*models.User, error)
	CreateUser(user *models.User) error
	UpdateUserBalance(username string, amount float64) error
	DeleteUser(username string) error
}

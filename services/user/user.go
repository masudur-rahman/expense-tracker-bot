package user

import (
	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/repos"
	"github.com/masudur-rahman/expense-tracker-bot/services"
)

type userService struct {
	userRepo repos.UserRepository
}

var _ services.UserService = &userService{}

func NewUserService(userRepo repos.UserRepository) *userService {
	return &userService{
		userRepo: userRepo,
	}
}

func (u *userService) GetUserByID(id string) (*models.User, error) {
	return u.userRepo.GetUserByID(id)
}

func (u *userService) GetUserByName(username string) (*models.User, error) {
	return u.userRepo.GetUserByName(username)
}

func (u *userService) ListUsers() ([]models.User, error) {
	return u.userRepo.ListUsers()
}

func (u *userService) CreateUser(user *models.User) error {
	return u.userRepo.AddNewUser(user)
}

func (u *userService) UpdateUserBalance(id string, amount float64) error {
	return u.userRepo.UpdateUserBalance(id, amount)
}

func (u *userService) DeleteUser(id string) error {
	return u.userRepo.DeleteUser(id)
}

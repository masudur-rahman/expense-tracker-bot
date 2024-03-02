package user

import (
	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/repos"
)

type userService struct {
	userRepo repos.UserRepository
}

func NewUserService(userRepo repos.UserRepository) *userService {
	return &userService{userRepo: userRepo}
}

func (us *userService) GetUserByID(id int64) (*models.User, error) {
	return us.userRepo.GetUserByID(id)
}

func (us *userService) GetUserByTelegramID(id int64) (*models.User, error) {
	filter := models.User{TelegramID: id}
	return us.userRepo.GetUser(filter)
}

func (us *userService) GetUserByUsername(username string) (*models.User, error) {
	filter := models.User{Username: username}
	return us.userRepo.GetUser(filter)
}

func (us *userService) ListUsers() ([]models.User, error) {
	return us.userRepo.ListUsers()
}

func (us *userService) SignUp(user *models.User) error {
	return us.userRepo.AddNewUser(user)
}

func (us *userService) UpdateUser(id int64, user *models.User) error {
	return us.userRepo.UpdateUser(id, user)
}

func (us *userService) DeleteUser(id int64) error {
	return us.userRepo.DeleteUser(id)
}

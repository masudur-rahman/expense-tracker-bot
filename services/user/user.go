package user

import (
	"strings"
	"time"

	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/models/gqtypes"
	"github.com/masudur-rahman/expense-tracker-bot/pkg"
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

func (us *userService) ValidateUser(params gqtypes.RegisterParams) error {
	_, err := us.userRepo.FindByName(params.Username)
	if err != nil && !models.IsErrNotFound(err) {
		return err
	} else if err == nil {
		return models.ErrUserAlreadyExist{Username: params.Username}
	}

	_, err = us.userRepo.FindByEmail(params.Email)
	if err != nil && !models.IsErrNotFound(err) {
		return err
	} else if err == nil {
		return models.ErrUserAlreadyExist{Email: params.Email}
	}

	return nil
}

func (us *userService) GetUser(id string) (*models.User, error) {
	user, err := us.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *userService) GetUserByName(username string) (*models.User, error) {
	user, err := us.userRepo.FindByName(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *userService) GetUserByEmail(email string) (*models.User, error) {
	user, err := us.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *userService) ListUsers(filter models.User, limit int64) ([]*models.User, error) {
	users, err := us.userRepo.FindUsers(filter)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (us *userService) CreateUser(params gqtypes.RegisterParams) (*models.User, error) {
	if err := us.ValidateUser(params); err != nil {
		return nil, err
	}

	user := &models.User{
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		Username:     params.Username,
		Email:        params.Email,
		PasswordHash: pkg.MustHashPassword(params.Password),
		IsActive:     true,
		IsAdmin:      false,
		CreatedUnix:  time.Now().Unix(),
	}
	if err := us.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *userService) UpdateUser(params gqtypes.UserParams) (*models.User, error) {
	user, err := us.userRepo.FindByName(params.Username)
	if err != nil {
		return nil, err
	}
	user.FirstName = params.FirstName
	user.LastName = params.LastName
	user.Location = params.Location
	user.Bio = params.Bio

	if err = us.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *userService) DeleteUser(id string) error {
	_, err := us.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	return us.userRepo.Delete(id)
}

func (us *userService) LoginUser(usernameOrEmail string, passwd string) (*models.User, error) {
	var err error
	var user *models.User
	if strings.Contains(usernameOrEmail, "@") {
		user, err = us.GetUserByEmail(usernameOrEmail)
		if err != nil {
			if models.IsErrNotFound(err) {
				return nil, models.ErrUserPasswordMismatch{}
			}
			return nil, err
		}
	} else {
		user, err = us.GetUserByName(usernameOrEmail)
		if err != nil {
			if models.IsErrNotFound(err) {
				return nil, models.ErrUserPasswordMismatch{}
			}
			return nil, err
		}
	}

	if !pkg.CheckPasswordHash(passwd, user.PasswordHash) {
		return nil, models.ErrUserPasswordMismatch{}
	}

	return user, nil
}

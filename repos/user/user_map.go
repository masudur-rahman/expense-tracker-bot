package user

import (
	"fmt"
	"net/http"

	"github.com/masudur-rahman/expense-tracker-bot/infra/database/mock"
	"github.com/masudur-rahman/expense-tracker-bot/infra/logr"
	"github.com/masudur-rahman/expense-tracker-bot/models"
)

type MapUserRepository struct {
	db     mock.Database
	logger logr.Logger
}

func NewMapUserRepository(db mock.Database, logger logr.Logger) *MapUserRepository {
	return &MapUserRepository{
		db:     db.Entity("user"),
		logger: logger,
	}
}

func (u *MapUserRepository) FindByID(id string) (*models.User, error) {
	u.logger.Infow("finding user by id", "id", id)
	var user models.User
	found, err := u.db.FindOne(&user)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, models.ErrUserNotFound{ID: id}
	}
	return &user, nil
}

func (u *MapUserRepository) FindByEmail(email string) (*models.User, error) {
	u.logger.Infow("finding user by email", "email", email)
	filter := models.User{
		Email: email,
	}
	var user models.User
	found, err := u.db.FindOne(&user, filter)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, models.ErrUserNotFound{Email: email}
	}
	return &user, nil
}

func (u *MapUserRepository) FindUsers(filter models.User) ([]*models.User, error) {
	u.logger.Infow("finding users by filter", "filter", fmt.Sprintf("%+v", filter))
	users := make([]*models.User, 0)
	err := u.db.FindMany(&users, filter)
	return users, err
}

func (u *MapUserRepository) Create(user *models.User) error {
	u.logger.Infow("creating user")
	_, err := u.db.InsertOne(user)
	return err
}

func (u *MapUserRepository) Update(user *models.User) error {
	u.logger.Infow("updating user")
	if user.ID == "" {
		return models.StatusError{
			Status:  http.StatusBadRequest,
			Message: "user id missing",
		}
	}

	return u.db.ID(user.ID).UpdateOne(user)
}

func (u *MapUserRepository) Delete(id string) error {
	u.logger.Infow("deleting user by id", "id", id)
	return u.db.ID(id).DeleteOne()
}

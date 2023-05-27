package user

import (
	"fmt"
	"net/http"

	isql "github.com/masudur-rahman/expense-tracker-bot/infra/database/sql"
	"github.com/masudur-rahman/expense-tracker-bot/infra/logr"
	"github.com/masudur-rahman/expense-tracker-bot/models"

	"github.com/rs/xid"
)

type SQLUserRepository struct {
	db     isql.Database
	logger logr.Logger
}

func NewSQLUserRepository(db isql.Database, logger logr.Logger) *SQLUserRepository {
	return &SQLUserRepository{
		db:     db.Table("user"),
		logger: logger,
	}
}

func (u *SQLUserRepository) FindByID(id string) (*models.User, error) {
	u.logger.Infow("finding user by id", "id", id)
	var user models.User
	found, err := u.db.ID(id).FindOne(&user)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, models.ErrUserNotFound{ID: id}
	}
	return &user, nil
}

func (u *SQLUserRepository) FindByName(username string) (*models.User, error) {
	u.logger.Infow("finding user by name", "name", username)
	filter := models.User{
		Username: username,
	}
	var user models.User
	found, err := u.db.FindOne(&user, filter)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, models.ErrUserNotFound{Username: username}
	}
	return &user, nil
}

func (u *SQLUserRepository) FindByEmail(email string) (*models.User, error) {
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

func (u *SQLUserRepository) FindUsers(filter models.User) ([]*models.User, error) {
	u.logger.Infow("finding users by filter", "filter", fmt.Sprintf("%+v", filter))
	users := make([]*models.User, 0)
	err := u.db.FindMany(&users, filter)
	return users, err
}

func (u *SQLUserRepository) Create(user *models.User) error {
	u.logger.Infow("creating user")
	if user.ID == "" {
		user.ID = xid.New().String()
	}
	user.XKey = user.ID
	id, err := u.db.InsertOne(user)
	u.logger.Infow("user created", "id", id)
	return err
}

func (u *SQLUserRepository) Update(user *models.User) error {
	u.logger.Infow("updating user")
	if user.ID == "" {
		return models.StatusError{
			Status:  http.StatusBadRequest,
			Message: "user id missing",
		}
	}

	return u.db.ID(user.ID).UpdateOne(user)
}

func (u *SQLUserRepository) Delete(id string) error {
	u.logger.Infow("deleting user by id", "id", id)
	return u.db.ID(id).DeleteOne()
}

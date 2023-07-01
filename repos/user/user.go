package user

import (
	"time"

	"github.com/masudur-rahman/expense-tracker-bot/infra/logr"
	"github.com/masudur-rahman/expense-tracker-bot/models"

	isql "github.com/masudur-rahman/database/sql"
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

func (u *SQLUserRepository) GetUserByID(id string) (*models.User, error) {
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

func (u *SQLUserRepository) GetUserByName(username string) (*models.User, error) {
	u.logger.Infow("finding user by name", "name", username)
	filter := models.User{
		Name: username,
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

func (u *SQLUserRepository) UpdateUserBalance(id string, txnAmount float64) error {
	u.logger.Infow("updating user")
	user, err := u.GetUserByID(id)
	if err != nil {
		return err
	}

	user.Balance += txnAmount
	user.LastTxnTimestamp = time.Now().Unix()

	return u.db.ID(user.ID).UpdateOne(user)
}

func (u *SQLUserRepository) AddNewUser(user *models.User) error {
	_, err := u.db.InsertOne(user)
	return err
}

func (u *SQLUserRepository) ListUsers() ([]models.User, error) {
	u.logger.Infow("listing users")
	users := make([]models.User, 0)
	err := u.db.FindMany(&users, models.User{})
	return users, err
}

func (u *SQLUserRepository) DeleteUser(id string) error {
	u.logger.Infow("deleting user", "id", id)
	return u.db.DeleteOne(models.User{ID: id})
}

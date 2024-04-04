package user

import (
	"github.com/masudur-rahman/expense-tracker-bot/infra/logr"
	"github.com/masudur-rahman/expense-tracker-bot/models"

	isql "github.com/masudur-rahman/styx/sql"
)

type SQLUserRepository struct {
	db     isql.Engine
	logger logr.Logger
}

func NewSQLUserRepository(db isql.Engine, logger logr.Logger) *SQLUserRepository {
	return &SQLUserRepository{
		db:     db.Table("user"),
		logger: logger,
	}
}

func (u *SQLUserRepository) GetUserByID(id int64) (*models.User, error) {
	u.logger.Infow("finding user by id", "id", id)
	var user models.User
	found, err := u.db.ID(id).FindOne(&user)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, models.ErrUserNotFound{}
	}
	return &user, nil
}

func (u *SQLUserRepository) GetUser(filter models.User) (*models.User, error) {
	//u.logger.Infow("finding user by telegram id", "telegram id", id)
	var user models.User
	found, err := u.db.FindOne(&user, filter)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, models.ErrUserNotFound{}
	}
	return &user, nil
}

func (u *SQLUserRepository) GetUserByUsername(username string) (*models.User, error) {
	u.logger.Infow("finding user by name", "username", username)
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

func (u *SQLUserRepository) ListUsers() ([]models.User, error) {
	u.logger.Infow("listing users")
	users := make([]models.User, 0)
	err := u.db.FindMany(&users)
	return users, err
}

func (u *SQLUserRepository) AddNewUser(user *models.User) error {
	_, err := u.db.InsertOne(user)
	return err
}

func (u *SQLUserRepository) UpdateUser(id int64, us *models.User) error {
	user, err := u.GetUserByID(id)
	if err != nil {
		return err
	}
	user.Username = us.Username
	user.FirstName = us.FirstName
	user.LastName = us.LastName

	return u.db.ID(id).UpdateOne(user)
}

func (u *SQLUserRepository) DeleteUser(id int64) error {
	u.logger.Infow("deleting user", "id", id)
	return u.db.DeleteOne(models.User{ID: id})
}

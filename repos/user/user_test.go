package user

import (
	"testing"

	"github.com/masudur-rahman/expense-tracker-bot/infra/database/nosql/mock"
	"github.com/masudur-rahman/expense-tracker-bot/infra/logr"
	"github.com/masudur-rahman/expense-tracker-bot/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func initializeDatabaseAndUserRepo(ctl *gomock.Controller) (*mock.MockDatabase, *NoSQLUserRepository) {
	db := mock.NewMockDatabase(ctl)
	db.EXPECT().Collection("user").Return(db).MaxTimes(2)
	ur := NewNoSQLUserRepository(db, logr.DefaultLogger)

	return db, ur
}

func TestNewNoSQLUserRepository(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db, ur := initializeDatabaseAndUserRepo(ctl)

	assert.NotNil(t, ur)
	assert.Equal(t, db.Collection("user"), ur.db)
	assert.Equal(t, logr.DefaultLogger, ur.logger)
}

func TestNoSQLUserRepository_FindByID(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db, ur := initializeDatabaseAndUserRepo(ctl)

	t.Run("user should not exist", func(t *testing.T) {
		id := "random-id"

		gomock.InOrder(
			db.EXPECT().ID(id).Return(db),
			db.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(false, models.ErrUserNotFound{ID: id}),
		)

		user, err := ur.FindByID(id)
		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrUserNotFound{ID: id})
		assert.Nil(t, user)
	})

	t.Run("user must exist", func(t *testing.T) {
		id := "abc-xyz"

		uf := models.User{
			ID:        id,
			FirstName: "Masudur",
			LastName:  "Rahman",
			Username:  "masud",
		}

		gomock.InOrder(
			db.EXPECT().ID(id).Return(db),
			db.EXPECT().FindOne(gomock.Any(), gomock.Any()).DoAndReturn(func(user *models.User, _ ...interface{}) (bool, error) {
				*user = uf
				return true, nil
			}),
		)

		user, err := ur.FindByID(id)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.EqualValues(t, &uf, user)
	})
}

func TestNoSQLUserRepository_FindByName(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db, ur := initializeDatabaseAndUserRepo(ctl)

	t.Run("user should not exist", func(t *testing.T) {
		name := "random-name"
		filter := models.User{Username: name}
		gomock.InOrder(
			db.EXPECT().FindOne(gomock.Any(), filter).Return(false, models.ErrUserNotFound{Username: name}),
		)

		user, err := ur.FindByName(name)
		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrUserNotFound{Username: name})
		assert.Nil(t, user)
	})

	t.Run("user must exist", func(t *testing.T) {
		name := "masud"

		filter := models.User{Username: name}
		uf := models.User{
			ID:        "abc-xyz",
			FirstName: "Masudur",
			LastName:  "Rahman",
			Username:  name,
		}

		gomock.InOrder(
			db.EXPECT().FindOne(gomock.Any(), filter).DoAndReturn(func(user *models.User, _ ...interface{}) (bool, error) {
				*user = uf
				return true, nil
			}),
		)

		user, err := ur.FindByName(name)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.EqualValues(t, &uf, user)
	})
}

func TestNoSQLUserRepository_FindByEmail(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db, ur := initializeDatabaseAndUserRepo(ctl)

	t.Run("user should not exist", func(t *testing.T) {
		email := "random@email.address"
		filter := models.User{Email: email}
		gomock.InOrder(
			db.EXPECT().FindOne(gomock.Any(), filter).Return(false, models.ErrUserNotFound{Email: email}),
		)

		user, err := ur.FindByEmail(email)
		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrUserNotFound{Email: email})
		assert.Nil(t, user)
	})

	t.Run("user must exist", func(t *testing.T) {
		email := "masudjuly02@gmail.com"

		filter := models.User{Email: email}
		uf := models.User{
			ID:        "abc-xyz",
			FirstName: "Masudur",
			LastName:  "Rahman",
			Username:  "masud",
			Email:     email,
		}

		gomock.InOrder(
			db.EXPECT().FindOne(gomock.Any(), filter).DoAndReturn(func(user *models.User, _ ...interface{}) (bool, error) {
				*user = uf
				return true, nil
			}),
		)

		user, err := ur.FindByEmail(email)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.EqualValues(t, &uf, user)
	})
}

func TestNoSQLUserRepository_FindUsers(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db, ur := initializeDatabaseAndUserRepo(ctl)

	t.Run("no users", func(t *testing.T) {
		filter := models.User{
			LastName: "NoName",
		}
		gomock.InOrder(
			db.EXPECT().FindMany(gomock.Any(), filter).Return(nil),
		)

		users, err := ur.FindUsers(filter)
		assert.NoError(t, err)
		assert.EqualValues(t, 0, len(users))
	})
}

func TestNoSQLUserRepository_Create(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db, ur := initializeDatabaseAndUserRepo(ctl)

	t.Run("should create user", func(t *testing.T) {
		id := "abc-xyz"
		user := models.User{
			ID:        id,
			FirstName: "Masudur",
			LastName:  "Rahman",
			IsActive:  true,
		}

		gomock.InOrder(
			db.EXPECT().InsertOne(gomock.Any()).Return(id, nil),
		)

		err := ur.Create(&user)
		assert.NoError(t, err)
		assert.Equal(t, id, user.XKey)
	})
}

func TestNoSQLUserRepository_Update(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db, ur := initializeDatabaseAndUserRepo(ctl)

	t.Run("should update user", func(t *testing.T) {
		id := "abc-xyz"
		user := models.User{
			ID:        id,
			FirstName: "Masudur",
			LastName:  "Rahman",
			IsActive:  true,
		}

		gomock.InOrder(
			db.EXPECT().ID(id).Return(db),
			db.EXPECT().UpdateOne(gomock.Any()).Return(nil),
		)

		err := ur.Update(&user)
		assert.NoError(t, err)
		assert.NotNil(t, user)
	})
}

func TestNoSQLUserRepository_Delete(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db, ur := initializeDatabaseAndUserRepo(ctl)

	t.Run("should delete a user", func(t *testing.T) {
		id := "abc-xyz"

		gomock.InOrder(
			db.EXPECT().ID(id).Return(db),
			db.EXPECT().DeleteOne().Return(nil),
		)

		err := ur.Delete(id)
		assert.NoError(t, err)
	})
}

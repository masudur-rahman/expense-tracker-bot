package user_test

import (
	"testing"

	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/models/gqtypes"
	"github.com/masudur-rahman/expense-tracker-bot/pkg"
	userRepo "github.com/masudur-rahman/expense-tracker-bot/repos/user"
	userSvc "github.com/masudur-rahman/expense-tracker-bot/services/user"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_userService_GetUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ur := userRepo.NewMockUserRepository(ctl)
	us := userSvc.NewUserService(ur)

	t.Run("user should not exist", func(t *testing.T) {
		id := "random-id"

		gomock.InOrder(
			ur.EXPECT().FindByID(id).Return(nil, models.ErrUserNotFound{ID: id}),
		)

		user, err := us.GetUser(id)
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
			ur.EXPECT().FindByID(id).Return(&uf, nil),
		)

		user, err := ur.FindByID(id)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.EqualValues(t, &uf, user)
	})
}

func Test_userService_GetUserByName(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ur := userRepo.NewMockUserRepository(ctl)
	us := userSvc.NewUserService(ur)

	t.Run("user should not exist", func(t *testing.T) {
		name := "random-name"

		gomock.InOrder(
			ur.EXPECT().FindByName(name).Return(nil, models.ErrUserNotFound{Username: name}),
		)

		user, err := us.GetUserByName(name)
		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrUserNotFound{Username: name})
		assert.Nil(t, user)
	})

	t.Run("user must exist", func(t *testing.T) {
		name := "masud"

		uf := models.User{
			ID:        "abc-xyz",
			FirstName: "Masudur",
			LastName:  "Rahman",
			Username:  name,
		}

		gomock.InOrder(
			ur.EXPECT().FindByName(name).Return(&uf, nil),
		)

		user, err := ur.FindByName(name)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.EqualValues(t, &uf, user)
	})
}

func Test_userService_GetUserByEmail(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ur := userRepo.NewMockUserRepository(ctl)
	us := userSvc.NewUserService(ur)

	t.Run("user should not exist", func(t *testing.T) {
		email := "nonexistent@example.com"

		gomock.InOrder(
			ur.EXPECT().FindByEmail(email).Return(nil, models.ErrUserNotFound{Email: email}),
		)

		user, err := us.GetUserByEmail(email)
		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrUserNotFound{Email: email})
		assert.Nil(t, user)
	})

	t.Run("user must exist", func(t *testing.T) {
		email := "masud@example.com"

		uf := models.User{
			ID:        "abc-xyz",
			FirstName: "Masudur",
			LastName:  "Rahman",
			Username:  "masud",
			Email:     email,
		}

		gomock.InOrder(
			ur.EXPECT().FindByEmail(email).Return(&uf, nil),
		)

		user, err := ur.FindByEmail(email)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.EqualValues(t, &uf, user)
	})
}

func Test_userService_ListUsers(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ur := userRepo.NewMockUserRepository(ctl)
	us := userSvc.NewUserService(ur)

	t.Run("no users found", func(t *testing.T) {
		var filter models.User
		limit := int64(10)

		gomock.InOrder(
			ur.EXPECT().FindUsers(filter).Return([]*models.User{}, nil),
		)

		users, err := us.ListUsers(filter, limit)
		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Len(t, users, 0)
	})

	t.Run("multiple users found", func(t *testing.T) {
		var filter models.User

		uf1 := models.User{
			ID:        "abc-xyz",
			FirstName: "Masudur",
			LastName:  "Rahman",
			Username:  "masud",
			Email:     "masud@example.com",
		}

		uf2 := models.User{
			ID:        "def-uvw",
			FirstName: "John",
			LastName:  "Doe",
			Username:  "johndoe",
			Email:     "johndoe@example.com",
		}

		gomock.InOrder(
			ur.EXPECT().FindUsers(filter).Return([]*models.User{&uf1, &uf2}, nil),
		)

		users, err := ur.FindUsers(filter)
		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Len(t, users, 2)
		assert.EqualValues(t, &uf1, users[0])
		assert.EqualValues(t, &uf2, users[1])
	})
}

func Test_userService_CreateUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	params := gqtypes.RegisterParams{
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
		Email:     "johndoe@example.com",
		Password:  "password",
	}

	ur := userRepo.NewMockUserRepository(ctl)
	us := userSvc.NewUserService(ur)

	t.Run("create user with valid parameters", func(t *testing.T) {
		gomock.InOrder(
			ur.EXPECT().FindByName(params.Username).Return(nil, models.ErrUserNotFound{Username: params.Username}),
			ur.EXPECT().FindByEmail(params.Email).Return(nil, models.ErrUserNotFound{Email: params.Email}),
			ur.EXPECT().Create(gomock.Any()).DoAndReturn(func(user *models.User) error {
				user.ID = "abc-xyz"
				return nil
			}),
		)

		user, err := us.CreateUser(params)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "abc-xyz", user.ID)
		assert.Equal(t, params.FirstName, user.FirstName)
		assert.Equal(t, params.LastName, user.LastName)
		assert.Equal(t, params.Username, user.Username)
		assert.Equal(t, params.Email, user.Email)
		assert.True(t, user.IsActive)
		assert.False(t, user.IsAdmin)
	})

	t.Run("create user with existing username", func(t *testing.T) {
		uf := models.User{
			ID:        "abc-xyz",
			FirstName: "Existing",
			LastName:  "User",
			Username:  params.Username,
			Email:     "existinguser@example.com",
		}

		gomock.InOrder(
			ur.EXPECT().FindByName(params.Username).Return(&uf, nil),
		)

		user, err := us.CreateUser(params)
		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrUserAlreadyExist{Username: params.Username})
		assert.Nil(t, user)
	})

	t.Run("create user with existing email", func(t *testing.T) {
		uf := models.User{
			ID:        "abc-xyz",
			FirstName: "Existing",
			LastName:  "User",
			Username:  "existinguser",
			Email:     params.Email,
		}

		gomock.InOrder(
			ur.EXPECT().FindByName(params.Username).Return(nil, models.ErrUserNotFound{Username: params.Username}),
			ur.EXPECT().FindByEmail(params.Email).Return(&uf, nil),
		)

		user, err := us.CreateUser(params)
		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrUserAlreadyExist{Email: params.Email})
		assert.Nil(t, user)
	})
}

func Test_userService_UpdateUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ur := userRepo.NewMockUserRepository(ctl)
	us := userSvc.NewUserService(ur)

	t.Run("update user location", func(t *testing.T) {
		username := "masud"
		existingUser := &models.User{
			Username:  username,
			FirstName: "Masudur",
			LastName:  "Rahman",
			Location:  "Chittagong",
		}

		updateParams := gqtypes.UserParams{
			Username:  username,
			FirstName: "Masudur",
			LastName:  "Rahman",
			Location:  "Dhaka",
		}

		gomock.InOrder(
			ur.EXPECT().FindByName(username).Return(existingUser, nil),
			ur.EXPECT().Update(gomock.Any()).Return(nil),
		)

		user, err := us.UpdateUser(updateParams)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.EqualValues(t, updateParams.Location, user.Location)
	})
}

func Test_userService_DeleteUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ur := userRepo.NewMockUserRepository(ctl)
	us := userSvc.NewUserService(ur)

	t.Run("delete user", func(t *testing.T) {
		userID := "abc-xyz"

		gomock.InOrder(
			ur.EXPECT().FindByID(userID).Return(&models.User{}, nil),
			ur.EXPECT().Delete(userID).Return(nil),
		)

		err := us.DeleteUser(userID)
		assert.NoError(t, err)
	})
}

func Test_userService_LoginUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ur := userRepo.NewMockUserRepository(ctl)
	us := userSvc.NewUserService(ur)

	t.Run("login user with email", func(t *testing.T) {
		usernameOrEmail := "masud@example.com"
		password := "secret"

		gomock.InOrder(
			ur.EXPECT().FindByEmail(usernameOrEmail).Return(&models.User{
				Email:        usernameOrEmail,
				PasswordHash: pkg.MustHashPassword(password),
				IsActive:     true,
			}, nil),
		)

		user, err := us.LoginUser(usernameOrEmail, password)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.EqualValues(t, usernameOrEmail, user.Email)
	})
}

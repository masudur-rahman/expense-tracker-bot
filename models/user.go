package models

import (
	"github.com/masudur-rahman/expense-tracker-bot/models/gqtypes"
)

type User struct {
	XKey      string `json:"_key"`
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Bio       string `json:"bio"`
	Location  string `json:"location"`
	Avatar    string `json:"avatar"`

	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`

	IsActive bool `json:"isActive"`
	IsAdmin  bool `json:"isAdmin"`

	CreatedUnix   int64 `json:"createdUnix"`
	UpdatedUnix   int64 `json:"updatedUnix"`
	LastLoginUnix int64 `json:"lastLoginUnix"`
}

func (u *User) APIFormat() gqtypes.User {
	return gqtypes.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		//FullName: fmt.Sprintf("%s %s", u.FirstName, u.LastName),
		Bio:      u.Bio,
		Location: u.Location,
		Avatar:   u.Avatar,
		IsActive: u.IsActive,
		IsAdmin:  u.IsAdmin,
	}
}

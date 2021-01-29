package model

import (
	"github.com/jmoiron/sqlx"
)

type (
	model struct {
		db *sqlx.DB
	}

	Database interface {
		SignupUser(user *User) (err error)
		LoginUser(user *User) (err error)
		UpdateUserProfile(user *User) (err error)
		UpdateUserPassword(user *User) (err error)
		UpdateUserLocale(user *User) (err error)
		GetUser(id uint64) (*User, error)
		GetDetailForUserLogin(email string) (*User, error)
	}
)

func NewModel(db *sqlx.DB) *model {
	return &model {
		db: db,
	}
}
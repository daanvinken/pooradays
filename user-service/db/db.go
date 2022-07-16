package database

import (
	"user-service/env"
	"user-service/model"
)

type Provider interface {
	Connect(e env.Provider) error
	CreateUser(u *model.User) (*model.User, error)
	GetUserById(id string) (*model.User, error)
	GetUserByName(id string) (*model.User, error)
	GetUserByEmail(id string) (*model.User, error)
	Signup(u *model.Signup) (*model.User, error)
	Login(u *model.Login) (*model.User, error)
}

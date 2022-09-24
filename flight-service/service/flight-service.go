package service

import (
	database "user-service/db"
	"user-service/pkg/model"
)

type UserService interface {
	Signup(s *model.Signup) (*model.User, error)
	Login(l *model.Login) (*model.UserAccess, error)
	GetUserById(id uint) (*model.User, error)
	UpdateUserById(id uint, column string, value interface{}) (*model.User, error)
}

/*
 *	User service layer to help interaction between user controller and databse.
**/
type userService struct {
}

var (
	Db database.Provider = database.NewPG()
)

func NewUserService() UserService {
	return &userService{}
}

func (s *userService) Signup(u *model.Signup) (*model.User, error) {
	return Db.Signup(u)

}

func (s *userService) Login(u *model.Login) (*model.UserAccess, error) {
	return Db.Login(u)
}

func (s *userService) GetUserById(id uint) (*model.User, error) {
	return Db.GetUserById(id)
}

func (s *userService) UpdateUserById(id uint, column string, value interface{}) (*model.User, error) {
	return Db.UpdateUserById(id, column, value)
}

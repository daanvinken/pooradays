package database

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"user-service/env"
	"user-service/pkg/model"
)

var instance *Postgres

/*
 *	User and Org business layer to accept request from service layer and persist in database.
**/

type Postgres struct {
	Db *gorm.DB
}

func NewPG() *Postgres {
	if instance != nil {
		return instance
	}
	instance = &Postgres{}
	return instance
}

func (pg *Postgres) autoMigrate() {
	pg.Db.AutoMigrate(&model.User{})
	fmt.Println("Automigrate complete")
}

func (pg *Postgres) Connect(e env.Provider) error {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Europe/Amsterdam",
		e.Get("DB.HOST"), e.Get("DB.USERNAME"), e.Get("DB.PASSWORD"), e.Get("DB.DATABASE"), e.Get("DB.PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
		return err
	}
	pg.Db = db
	_ = db
	fmt.Println("Connected to db")
	instance.autoMigrate()
	return nil
}

func (pg *Postgres) CreateUser(user *model.User) (*model.User, error) {
	if e := pg.Db.Create(&user).Error; e != nil {
		return nil, e
	}

	return user, nil
}

func (pg *Postgres) GetUserById(id uint) (*model.User, error) {
	var user *model.User
	if e := pg.Db.First(&user, "id = ?", id).Error; e != nil {
		return nil, e
	}
	return user, nil
}

func (pg *Postgres) GetUserByToken(token string) (*model.User, error) {
	var user *model.User
	if e := pg.Db.First(&user, "token = ?", token).Error; e != nil {
		return nil, e
	}
	return user, nil
}

func (pg *Postgres) UpdateUserById(id uint, column string, value interface{}) (*model.User, error) {
	var user *model.User
	var err error
	user, err = pg.GetUserById(id)
	if err != nil {
		return nil, errors.New("Could not find user.")
	}
	if e := pg.Db.Model(&user).Update(column, value).Error; e != nil {
		return nil, e
	}
	return user, nil
}

//func (pg *Postgres) GetUserByName(name string) (*model.User, error) {
//	var user *model.User
//	if e := pg.Db.Where(&model.User{Name: name}).First(&user).Error; e != nil {
//		return nil, e
//	}
//	return user, nil
//}

func (pg *Postgres) GetUserByEmail(email string) (*model.User, error) {
	var user *model.User
	if e := pg.Db.Where(&model.User{Email: email}).First(&user).Error; e != nil {
		return nil, e
	}
	return user, nil
}

func (pg *Postgres) Signup(u *model.Signup) (*model.User, error) {
	var e error
	var user *model.User

	tx := pg.Db.Begin()
	if user, e = pg.CreateUser(&u.User); e != nil {
		tx.Rollback()
		return user, e
	}

	tx.Commit()
	return user, nil
}

func (pg *Postgres) Login(u *model.Login) (*model.UserAccess, error) {
	var user *model.User
	var err error
	var userAccess *model.UserAccess

	if u.Password == "" {
		return nil, errors.New("Password is required")
	}

	if u.Email != "" {
		user, err = pg.GetUserByEmail(u.Email)
	} else {
		return nil, errors.New("Email is required")
	}

	if err != nil {
		return nil, errors.New("User doesn't exists")
	}
	if !model.ComparePwd(user.Password, []byte(u.Password)) {
		return nil, errors.New("Invalid password!")
	}
	userAccess, err = user.GetUserAccessToken()

	if err == nil {

	}
	return userAccess, err

}

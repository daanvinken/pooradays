package database

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"user-service/env"
	"user-service/model"
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
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Kolkata",
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

func (pg *Postgres) GetUserById(id string) (*model.User, error) {
	var user *model.User
	if e := pg.Db.First(&user, "id = ?", id).Error; e != nil {
		return nil, e
	}
	user.Password = ""
	return user, nil
}

func (pg *Postgres) GetUserByName(name string) (*model.User, error) {
	var user *model.User
	if e := pg.Db.Where(&model.User{Name: name}).First(&user).Error; e != nil {
		return nil, e
	}
	return user, nil
}

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
		user.Password = ""
		return user, e
	}

	tx.Commit()
	user.Password = ""
	return user, nil
}

func (pg *Postgres) Login(u *model.Login) (*model.User, error) {
	var user *model.User
	var err error

	if u.Password == "" {
		return nil, errors.New("password is required")
	}

	if u.Email != "" {
		user, err = pg.GetUserByEmail(u.Email)
	} else if u.Name != "" {
		user, err = pg.GetUserByName(u.Name)
	} else {
		return nil, errors.New("name or email is required")
	}

	if err != nil {
		return nil, errors.New("user doesn't exists")
	}
	if model.ComparePwd(user.Password, []byte(u.Password)) {
		user.Password = ""
		return user, nil
	}
	return nil, errors.New("invalid password")
}

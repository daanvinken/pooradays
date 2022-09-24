package model

import (
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	FirstNames           string `json:"FirstName"`
	LastName             string
	Email                string `gorm:"unique" json:"email"`
	NewsletterSubscribed bool
	Password             string `json:"password,omitempty"`
	PrimaryCity          string
	Token                uuid.UUID
	TokenExpiryEpoch     int64
	LastLoginEpoch       int64 `gorm:"autoUpdateTime:milli"`
	Birthday             time.Time
	PassportNumber       string
	Nationality          string
	gorm.Model
}

type UserAccess struct {
	Id          uint
	Token       uuid.UUID
	TokenExpiry int64
}

func (u *User) GetUserAccessToken() (*UserAccess, error) {
	var out *UserAccess
	token, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	expiryDate := time.Now().Local().Add(time.Hour * time.Duration(1))
	userAccess := UserAccess{
		Id:          u.ID,
		Token:       token,
		TokenExpiry: expiryDate.Unix(),
	}
	out = &userAccess
	return out, nil
}

type Signup struct {
	User User
}

type Login struct {
	ID       int64
	Email    string
	Password string
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	token, _ := uuid.NewV4()
	u.Token = token
	u.Password = HashPwd([]byte(u.Password))
	return
}

func HashPwd(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func ComparePwd(hash string, pwd []byte) bool {
	byteHash := []byte(hash)

	err := bcrypt.CompareHashAndPassword(byteHash, pwd)
	return err == nil
}

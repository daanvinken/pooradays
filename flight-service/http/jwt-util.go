package http

import (
	"flight-service/client"
	jwt "github.com/dgrijalva/jwt-go"
)

func GetToken(name string) (string, error) {
	signingKey := []byte("keymaker")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": name,
		"role": "redpill",
	})
	tokenString, err := token.SignedString(signingKey)
	return tokenString, err
}

func VerifyToken(tokenString string) (bool, error) {
	//signingKey := []byte("keymaker")
	//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	//	return signingKey, nil
	//})
	var is_valid bool
	var err error
	if is_valid, err = client.VerifyToken(tokenString); err != nil {
		return false, err
	}
	return is_valid, nil

}

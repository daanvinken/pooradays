package client

import "net/http"

func VerifyToken(token string) (bool, error) {
	http.Get("localhost:8080/")
	return true, nil
}

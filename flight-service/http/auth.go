package http

import (
	"net/http"
	"strings"
)

func authenticate(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("programName")
	password := r.FormValue("programPassword")

	if len(name) == 0 || len(password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please provide name and password to obtain the token"))
		return
	}
	if (name == "neo" && password == "keanu") || (name == "morpheus" && password == "lawrence") {
		token, err := GetToken(name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error generating JWT token: " + err.Error()))
		} else {
			w.Header().Set("Authorization", "Bearer "+token)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Token: " + token))
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Name and password do not match"))
		return
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authorization Header"))
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		is_valid, err := VerifyToken(tokenString)
		if err != nil || is_valid == false {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error verifying JWT token: " + err.Error()))
			return
		}
		//name := claims.(jwt.MapClaims)["name"].(string)
		//role := claims.(jwt.MapClaims)["role"].(string)
		//
		//r.Header.Set("name", name)
		//r.Header.Set("role", role)

		next.ServeHTTP(w, r)
	})
}

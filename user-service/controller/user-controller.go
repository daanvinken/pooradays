package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"user-service/pkg/model"
	"user-service/service"
)

/*
 *	User controller layer to accept request from exposed API and pass it user service layer
**/

var (
	userSVC service.UserService = service.NewUserService()
)

func Health(w http.ResponseWriter, r *http.Request) {
	RespondWithStatus(w, http.StatusOK, "Server OK")
}

func VerifyToken(w http.ResponseWriter, r *http.Request) {
	var err error
	var is_valid bool
	token := r.URL.Query().Get("token")
	log.Println(token)

	if is_valid, err = userSVC.VerifyToken(token); err != nil {
		if err.Error() == "record not found" {
			goto ESCAPE_DB_ERROR
		}
		RespondWithError(w, http.StatusInternalServerError, err)
		return
	}
ESCAPE_DB_ERROR:
	if is_valid {
		RespondWithStatus(w, http.StatusOK, "Token is valid.")
	}
	RespondWithStatus(w, http.StatusUnauthorized, "Token is not valid.")

}

func Signup(w http.ResponseWriter, r *http.Request) {
	var u model.Signup
	var err error

	if err = json.NewDecoder(r.Body).Decode(&u); err != nil {
		fmt.Println("err ", err)
		RespondWithError(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()
	if _, err = userSVC.Signup(&u); err != nil {
		log.Println("Error during signup of user: %v", err)
		RespondWithError(w, http.StatusConflict, err)
		return
	}
	RespondWithStatus(w, http.StatusOK, "Great success!")
}

//func Login(w http.ResponseWriter, r *http.Request) {
//	var u model.Login
//	var err error
//	var user *model.User
//
//	if err = json.NewDecoder(r.Body).Decode(&u); err != nil {
//		RespondWithError(w, http.StatusBadRequest, err)
//		return
//	}
//	defer r.Body.Close()
//	if user, err = userSVC.Login(&u); err != nil {
//		RespondWithError(w, http.StatusInternalServerError, err)
//		return
//	}
//	RespondWithJSON(w, http.StatusOK, user)
//}

func Login(w http.ResponseWriter, r *http.Request) {
	var u model.Login
	var err error
	var userAccess *model.UserAccess

	if err = json.NewDecoder(r.Body).Decode(&u); err != nil {
		RespondWithError(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()
	if userAccess, err = userSVC.Login(&u); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err)
		return
	}
	if _, err = userSVC.UpdateUserById(userAccess.Id, "Token", userAccess.Token); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, userAccess)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	//var err error
	//var userId uint
	id := r.URL.Query().Get("id")
	userId, err := strconv.ParseUint(id, 0, 64)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err)
	}

	var user *model.User
	if user, err = userSVC.GetUserById(uint(userId)); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err)
		return
	}
	RespondWithJSON(w, http.StatusOK, user)

}

func RespondWithError(w http.ResponseWriter, code int, err error) {
	var error = ConvertErrorMessage(err)
	log.Println(error)
	var message = error.Error()
	RespondWithJSON(w, code, map[string]string{"message": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithStatus(w http.ResponseWriter, code int, status string) {
	RespondWithJSON(w, code, map[string]string{"message": status})
}

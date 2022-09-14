package controller

import (
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

/*
 *	Error converter to hide raw error handling away from api, but sometimes still provide a descriptive message
**/

func IsPgErrorCode(err error, errorCode string) bool {
	if pgerr, ok := err.(*pgconn.PgError); ok {
		return pgerr.Code == errorCode
	}
	return false
}

func ConvertErrorMessage(err error) error {
	switch {
	case IsPgErrorCode(err, pgerrcode.UniqueViolation):
		return errors.New("Error: User already exists.")
	default:
		return errors.New(fmt.Sprintf("An internal error occured: %v", err.Error()))
	}
}

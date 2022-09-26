package client

import (
	"log"
	"net/http"
	"os"
)

func VerifyToken(token string) (bool, error) {
	req, err := http.NewRequest("GET", "http://pooradays-user:8080/user/verify_token", nil)
	if err != nil {
		log.Print(err)
		return false, err
	}
	log.Println("Received request on " + os.Getenv("POD_NAME"))
	
	q := req.URL.Query()
	q.Add("token", token)
	req.URL.RawQuery = q.Encode()
	resp, err := http.Get(req.URL.String())
	if resp.StatusCode == http.StatusOK {
		return true, err
	}
	return false, err

}

package controller

import (
	"encoding/json"
	"flight-service/service"
	"log"
	"net/http"
	"os"
)

/*
 *	Flight controller layer to accept request from exposed API and pass it user service layer
**/

var (
	flightSVC service.FlightService = service.NewFlightService()
)

func Health(w http.ResponseWriter, r *http.Request) {
	RespondWithStatus(w, http.StatusOK, "Server OK")
}

func GetFlights(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request on " + os.Getenv("POD_NAME"))
	var err error
	if _, err = flightSVC.FindFlights("hey"); err != nil {
		log.Println("Error during flight search: %v", err)
		RespondWithError(w, http.StatusConflict, err)
		return
	}
	RespondWithStatus(w, http.StatusOK, "Great success!")
}

func RespondWithError(w http.ResponseWriter, code int, err error) {
	RespondWithJSON(w, code, map[string]string{"message": err.Error()})
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

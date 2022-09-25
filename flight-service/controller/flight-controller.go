package controller

import (
	"encoding/json"
	"flight-service/service"
	"log"
	"net/http"
	"time"
)

/*
 *	Flight controller layer to accept request from exposed API and pass it user service layer
**/

var (
	flightSVC service.FlightService = service.NewFlightService()
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	log.Println("Received health check at " + string(time.Now().Local().Unix()))
	RespondWithStatus(w, 200, "OK")
}

func GetFlights(w http.ResponseWriter, r *http.Request) {
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

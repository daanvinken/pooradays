package service

type FlightService interface {
	FindFlights(value interface{}) (string, error)
}

/*
 *	Flight service layer to help interaction between user controller and external APi
**/
type flightService struct {
}

func NewFlightService() FlightService {
	return &flightService{}
}

func (s *flightService) FindFlights(value interface{}) (string, error) {
	return "Got you a flight", nil
}

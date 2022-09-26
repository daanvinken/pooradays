package main

import (
	"flight-service/controller"
	"flight-service/env"
	"flight-service/http"
	"flight-service/service"
)

var (
	en         env.Provider = env.NewEnv()
	mainRouter              = http.NewMuxRouter()
)

func main() {
	initApp()
	initRoutes()

	mainRouter.Serve()
}

func initApp() {
	en.Init()
	service.NewFlightService()
}

func initRoutes() {
	mainRouter.UnsafeGet("/health", controller.Health)
	flightRouter := mainRouter.RegisterSubRoute("/flight")
	flightRouter.Get("/search", controller.GetFlights)
}

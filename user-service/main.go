package main

import (
	"user-service/controller"
	database "user-service/db"
	"user-service/env"
	"user-service/http"
	"user-service/service"
)

var (
	en         env.Provider      = env.NewEnv()
	db         database.Provider = database.NewPG()
	mainRouter                   = http.NewMuxRouter()
)

func main() {
	initApp()
	initRoutes()

	mainRouter.Serve()
}

func initApp() {
	en.Init()
	db.Connect(en)
	service.NewUserService()
}

func initRoutes() {
	mainRouter.Get("/health", controller.Health)

	userRouter := mainRouter.RegisterSubRoute("/user")
	userRouter.Post("/signup", controller.Signup)
	userRouter.Post("/login", controller.Login)
	userRouter.Get("/{id}", controller.GetUserById)
}

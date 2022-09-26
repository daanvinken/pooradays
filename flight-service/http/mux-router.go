package http

import (
	"fmt"
	"log"
	"net/http"

	"flight-service/env"
	"github.com/gorilla/mux"
)

type MuxRouter struct {
	router *mux.Router
}

var (
	Env env.Provider = env.NewEnv()
)

func NewMuxRouter() Router {
	return &MuxRouter{router: mux.NewRouter()}
}

func (mx *MuxRouter) UnsafeGet(url string, f func(w http.ResponseWriter, r *http.Request)) {
	mx.router.HandleFunc(url, f).Methods("GET")
}

func (mx *MuxRouter) Get(url string, f func(w http.ResponseWriter, r *http.Request)) {
	mx.router.Handle(url, authMiddleware(http.HandlerFunc(f))).Methods("GET")
}

func (mx *MuxRouter) Post(url string, f func(w http.ResponseWriter, r *http.Request)) {
	mx.router.Handle(url, authMiddleware(http.HandlerFunc(f))).Methods("POST")
}

func (mx *MuxRouter) RegisterSubRoute(path string) Router {
	subRouter := mx.router.PathPrefix(path).Subrouter()
	return &MuxRouter{router: subRouter}
}

func (mx *MuxRouter) Serve() {
	fmt.Printf("Server starting on %v:%v", Env.Get("HOST"), Env.Get("PORT"))
	fmt.Println()
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", Env.Get("HOST"), Env.Get("PORT")), mx.router))
}

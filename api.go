package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type App struct {
	Router *mux.Router
}

func initialize(a *App) {

	a.Router = mux.NewRouter()
	initializeRoutes(a)

}

func initializeRoutes(a *App) {
	a.Router.HandleFunc("/", standardHandler)
	a.Router.HandleFunc("/api/health", healthHandler).Methods("GET")

}

func standardHandler(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(TestMessage{
		Message: "Hello, World!",
	})
	return
}

func healthHandler(_ http.ResponseWriter, _ *http.Request) {

	return
}

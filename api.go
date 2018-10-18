package main

import (
	"encoding/json"
	"net/http"

	"github.com/real-time-footfall-analysis/rtfa-backend/locationupdate"

	"github.com/gorilla/mux"
	"github.com/real-time-footfall-analysis/rtfa-backend/eventstaticdata"
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
	eventstaticdata.Init(a.Router)
	locationupdate.Init(a.Router)
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

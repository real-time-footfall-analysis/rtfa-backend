package main

import (
	"encoding/json"
	"github.com/real-time-footfall-analysis/rtfa-backend/emergency"
	"github.com/real-time-footfall-analysis/rtfa-backend/notifications"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/real-time-footfall-analysis/rtfa-backend/eventlivedata"
	"github.com/real-time-footfall-analysis/rtfa-backend/eventstaticdata"
	"github.com/real-time-footfall-analysis/rtfa-backend/locationupdate"
	"github.com/real-time-footfall-analysis/rtfa-backend/readanalytics"
	"github.com/real-time-footfall-analysis/rtfa-backend/utils"
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
	a.Router.Methods("OPTIONS").HandlerFunc(preflightHandler)
	eventstaticdata.Init(a.Router)
	locationupdate.Init(a.Router)
	eventlivedata.Init(a.Router)
	readanalytics.Init(a.Router)
	emergency.Init(a.Router)
	notifications.Init(a.Router)
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {

	utils.SetAccessControlHeaders(w)

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

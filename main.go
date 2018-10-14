package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	eventmetadata "github.com/real-time-footfall-analysis/rtfa-backend/event-metadata"
)

type TestMessage struct {
	Message string `json:"message"`
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", standardHandler)
	r.HandleFunc("/api/health", healthHandler).Methods("GET")
	eventmetadata.Init(r)

	log.Fatal(http.ListenAndServe(":80", r))
}

func standardHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(TestMessage{
		Message: "Hello, World!",
	})
	return
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	return
}

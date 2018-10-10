package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type TestMessage struct {
	Message string `json:"message"`
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", standardHandler)
	r.HandleFunc("/api/health", healthHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
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

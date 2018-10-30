package readanalytics

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func Init(r *mux.Router) {

	r.HandleFunc("/events/{eventId}/tasks/{taskId}", getTaskResultHandler).Methods("GET")

}

func getTaskResultHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	eventIDStr := vars["eventId"]
	taskIDStr := vars["taskId"]

	if eventIDStr == "" {
		http.Error(
			w,
			fmt.Sprint("Missing event ID"),
			http.StatusBadRequest)
		return
	}

	if taskIDStr == "" {
		http.Error(
			w,
			fmt.Sprint("Missing task ID"),
			http.StatusBadRequest)
		return
	}

	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Failed to parse event id: %s", err),
			http.StatusBadRequest)
		return
	}

	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Failed to parse task id: %s", err),
			http.StatusBadRequest)
		return
	}

	result, err := fetchAnalyticsResult(eventID, taskID)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Failed to get region by event ID: %s", err),
			http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)

}

package readanalytics

import (
	"encoding/json"
	"fmt"
	"github.com/real-time-footfall-analysis/rtfa-backend/dynamoDB"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/real-time-footfall-analysis/rtfa-backend/utils"
)

var analytics_database dynamoDB.DynamoDBInterface = &dynamoDB.DynamoDBClient{}

func Init(r *mux.Router) {
	_ = analytics_database.InitConn("analytics_results")
	r.HandleFunc("/events/{eventId}/tasks/{taskId}", getTaskResultHandler).Methods("GET")
}

func getTaskResultHandler(w http.ResponseWriter, r *http.Request) {
	// Allow cross origin
	utils.SetAccessControlHeaders(w)

	// Get the task identifiers
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

	// Get the result
	pKeyColName := "EventID-TaskID"
	pKeyValue := fmt.Sprintf("%d-%d", eventID, taskID)
	result := analytics_database.GetItem(pKeyColName, pKeyValue)

	// Rename the results to a more usable format
	delete(result, pKeyColName)
	result["eventID"] = eventID
	result["taskID"] = taskID

	// Send the result back
	_ = json.NewEncoder(w).Encode(result)
}

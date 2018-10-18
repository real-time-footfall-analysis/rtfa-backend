package eventlivedata

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var db liveDbAdapter = &dynamoDbAdaptor{}

func Init(r *mux.Router) {

	db.initConn()

	r.HandleFunc("/live/heatmap/{eventId}", heatmapHandler).Methods("GET")
}

type mapRequest struct {
	EventID int32 `json:"event_id,string"`
}

func heatmapHandler(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["eventId"])
	if err != nil {
		log.Println("Cannot decode request event id:", err)
		http.Error(
			writer,
			fmt.Sprintf("Failed to decode request: %s", err),
			http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Println("Cannot decode request:", err)
		http.Error(
			writer,
			fmt.Sprintf("Failed to decode request: %s", err),
			http.StatusBadRequest)
		return
	}

	heatMap, _ := db.getLiveHeatMap(id)

	json.NewEncoder(writer).Encode(heatMap)

	return
}

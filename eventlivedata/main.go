package eventlivedata

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var db liveDbAdapter = &dynamoDbAdaptor{}

func Init(r *mux.Router) {

	db.initConn()

	r.HandleFunc("/live/heatmap", heatmapHandler).Methods("GET")
}

type mapRequest struct {
	EventID int32 `json:"event_id,string"`
}

func heatmapHandler(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	var req mapRequest

	err := decoder.Decode(&req)

	if err != nil {
		log.Println("Cannot decode request:", err)
		http.Error(
			writer,
			fmt.Sprintf("Failed to decode request: %s", err),
			http.StatusBadRequest)
		return
	}

	heatMap, _ := db.getLiveHeatMap(int(req.EventID))

	json.NewEncoder(writer).Encode(heatMap)

	return
}

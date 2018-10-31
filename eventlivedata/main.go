package eventlivedata

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/real-time-footfall-analysis/rtfa-backend/utils"
)

var db liveDbAdapter = &dynamoDbAdaptor{}

func Init(r *mux.Router) {

	db.initConn()

	r.HandleFunc("/live/heatmap/{eventId}", heatmapHandler).Methods("GET")
}

func heatmapHandler(writer http.ResponseWriter, request *http.Request) {

	utils.SetAccessControlHeaders(writer)

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

	heatMap, _ := db.getLiveHeatMap(id)

	json.NewEncoder(writer).Encode(heatMap)

	return
}

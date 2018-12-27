package eventlivedata

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/real-time-footfall-analysis/rtfa-backend/dynamoDB"
	"github.com/real-time-footfall-analysis/rtfa-backend/locationupdate"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/real-time-footfall-analysis/rtfa-backend/utils"
)

var db dynamoDB.DynamoDBInterface = &dynamoDB.DynamoDBClient{}

func Init(r *mux.Router) {
	// Create a connection to the database
	err := db.InitConn("current_position")
	if err != nil {
		log.Println("Error connecting to current position table")
		os.Exit(1)
	}

	// Register the endpoint
	r.HandleFunc("/live/heatmap/{eventId}", heatmapHandler).Methods("GET")
}

func heatmapHandler(writer http.ResponseWriter, request *http.Request) {
	// Allow cross origin
	utils.SetAccessControlHeaders(writer)

	// Decode the request
	vars := mux.Vars(request)
	eventId, err := strconv.Atoi(vars["eventId"])
	if err != nil {
		log.Println("Cannot decode request event id:", err)
		http.Error(
			writer,
			fmt.Sprintf("Failed to decode request: %s", err),
			http.StatusBadRequest)
		return
	}

	// Get the whole table
	unparsedRows := db.GetTableScan()
	var parsedRows []locationupdate.Movement_update = make([]locationupdate.Movement_update, len(unparsedRows))
	for index, row := range unparsedRows {
		_ = mapstructure.Decode(row, &parsedRows[index])
	}

	// Count the rows to make a heatmap
	regionCounts := make(map[int]int, 0)
	for _, row := range parsedRows {
		if eventId == *row.EventID {
			regionId := *row.RegionID
			count, ok := regionCounts[regionId]
			if !ok {
				regionCounts[regionId] = 1
			} else {
				regionCounts[regionId] = count + 1
			}
		}
	}

	// Return the result
	_ = json.NewEncoder(writer).Encode(regionCounts)
}

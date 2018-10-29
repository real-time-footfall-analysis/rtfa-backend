package locationupdate

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Init registers the endpoints exposed by this package
// with the given Router.
// Also initialises the static data database connection
var queue queue_adapter = &kenisis_queue{}

func Init(r *mux.Router) {

	queue.initConn()

	r.HandleFunc("/update", updateHandler).Methods("POST")
}

type update struct {
	UUID       string `json:"uuid"`
	EventID    int    `json:"eventId"`
	RegionID   int    `json:"regionId"`
	Entering   bool   `json:"entering"`
	OccurredAt int    `json:"occurredAt"`
}

func updateHandler(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	var update update

	err := decoder.Decode(&update)

	if err != nil {
		log.Println("Cannot decode update:", err)
		http.Error(
			writer,
			fmt.Sprintf("Failed to decode update: %s", err),
			http.StatusBadRequest)
		return
	}

	// TODO: replace with actual timestamp from frontend
	update.OccurredAt = int(time.Now().Unix())

	queue.addLocationUpdate(&update)

	return
}

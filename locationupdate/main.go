package locationupdate

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
	UUID     string `json:"uuid"`
	EventID  int    `json:"eventId"`
	RegionID int    `json:"regionId"`
	Entering bool   `json:"entering"`
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

	queue.addLocationUpdate(&update)

	return
}

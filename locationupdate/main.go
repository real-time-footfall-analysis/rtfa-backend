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
func Init(r *mux.Router) {

	initConn()

	r.HandleFunc("/update", updateHandler).Methods("POST")
}

func updateHandler(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	var update BluetoothUpdate

	err := decoder.Decode(&update)

	if err != nil {
		log.Println("Cannot decode update:", err)
		http.Error(
			writer,
			fmt.Sprintf("Failed to decode update: %s", err),
			http.StatusBadRequest)
		return
	}

	//TODO: check if event id is valid?

	addLocationUpdate(&update)

	return
}

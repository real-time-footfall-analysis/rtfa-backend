package eventmetadata

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Init registers the endpoints exposed by this package
// with the given Router.
func Init(r *mux.Router) {
	r.HandleFunc("/events", getEventsHandler).Methods("GET")
	r.HandleFunc("/events", postEventsHandler).Methods("POST")
	r.HandleFunc("/events/{eventId}", getEventHandler).Methods("GET")
	r.HandleFunc("/events/{eventId}/map", postEventMapHandler).Methods("POST")
	r.HandleFunc("/events/{eventId}/map", getEventMapHandler).Methods("GET")
}

func getEventsHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: handle

}

func postEventsHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: handle

}

func getEventHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: handle

}

func postEventMapHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: handle

}

func getEventMapHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: handle

}

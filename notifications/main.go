package notifications

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"github.com/real-time-footfall-analysis/rtfa-backend/dynamoDB"
	"github.com/real-time-footfall-analysis/rtfa-backend/pusher"
	"github.com/real-time-footfall-analysis/rtfa-backend/utils"
	"hash/fnv"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

type organiser_notification struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	RegionIds      []int  `json:"regionIds"`
	OccurredAt     int    `json:"occurredAt"`
	NotificationId int    `json:"notificationId"`
	EventId        int    `json:"eventId"`
}

var db dynamoDB.DynamoDBInterface = &dynamoDB.DynamoDBClient{}
var pb pusher.PusherBeamsInterface = &pusher.PusherBeamsClient{}
var pc pusher.PusherChannelInterface = &pusher.PusherChannelClient{}

func Init(r *mux.Router) {
	err := db.InitConn("notifications")
	if err != nil {
		os.Exit(1)
	}

	pb.InitConn()
	pc.InitConn()
	r.HandleFunc("/events/{eventId}/notifications", postNotification).Methods("POST")
	r.HandleFunc("/events/{eventId}/notifications", getAllNotifications).Methods("GET")
}

func postNotification(writer http.ResponseWriter, request *http.Request) {
	// Allow cross origin
	utils.SetAccessControlHeaders(writer)

	notification, err := decodeNotification(writer, request)
	if err != nil {
		return
	}

	// Check remaining fields and send response
	err = validateNotification(notification, writer)
	if err != nil {
		return
	}

	// Get the event id
	vars := mux.Vars(request)
	eventId, err := parseRequestArgs(vars, "eventId", writer)
	if err != nil {
		return
	}
	notification.EventId = eventId

	// Send the notification to pusher beams
	regions := intsToStrings(notification.RegionIds)
	publishId, err := pb.SendNotification(regions, notification.Title, notification.Description)

	// Generate a hash based on the response
	notification.NotificationId = hashString(publishId)

	// Send the item to the database
	db.SendItem(notification)

	// Send the notification to web app through pusher
	channelName := strconv.Itoa(eventId)
	data, _ := json.Marshal(notification)
	pc.SendItem(channelName, "organiser-notification", data)

	// Return the update to the user
	_ = json.NewEncoder(writer).Encode(notification)
}

func decodeNotification(writer http.ResponseWriter, request *http.Request) (organiser_notification, error) {
	// Try and decode the data
	var notification organiser_notification
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&notification)
	if err != nil {
		log.Println("Cannot decode organiser notification:", err)
		http.Error(
			writer,
			fmt.Sprintf("Failed to decode organiser notification: %s", err),
			http.StatusBadRequest)
	}
	return notification, err
}

func hashString(str string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(str))
	return int(h.Sum32())
}

func intsToStrings(ints []int) []string {
	var strings []string = make([]string, len(ints))
	for i := 0; i < len(ints); i++ {
		strings[i] = strconv.Itoa(ints[i])
	}
	return strings
}

func validateNotification(notification organiser_notification, writer http.ResponseWriter) error {
	// Check the fields in the data
	msg := ""
	if len(notification.Title) == 0 {
		msg = "title is empty"
	} else if len(notification.Description) == 0 {
		msg = "Description is empty"
	} else if notification.RegionIds == nil || len(notification.RegionIds) == 0 {
		msg = "No regions specified"
	} else if notification.OccurredAt == 0 {
		msg = "occurredAt timestamp missing"
	}

	// Send the error message back
	if len(msg) != 0 {
		log.Println(msg)
		http.Error(
			writer,
			fmt.Sprintf(msg),
			http.StatusBadRequest)
		return errors.New(msg)
	}
	return nil
}

func getAllNotifications(writer http.ResponseWriter, request *http.Request) {

	// Allow cross origin
	utils.SetAccessControlHeaders(writer)

	// Get the event id
	vars := mux.Vars(request)
	eventId, err := parseRequestArgs(vars, "eventId", writer)
	if err != nil {
		return
	}

	// Scan the table, parse the result and filter out old events
	unparsedRows := db.GetTableScan()
	var parsedRows []organiser_notification = make([]organiser_notification, len(unparsedRows))
	for index, row := range unparsedRows {
		_ = mapstructure.Decode(row, &parsedRows[index])
	}

	// Extract relevant values
	eventNotifications := extractRecentUpdates(parsedRows, eventId)

	// Sort the values by timestamp descending
	sort.Slice(eventNotifications, func(i, j int) bool {
		return eventNotifications[i].OccurredAt > eventNotifications[j].OccurredAt
	})

	// Transmit the result back
	_ = json.NewEncoder(writer).Encode(eventNotifications)
}

func extractRecentUpdates(parsed []organiser_notification, event int) (res []organiser_notification) {
	// Remove the values that don't satisfy a criteria
	deleted := 0
	for i := range parsed {
		j := i - deleted
		if parsed[j].EventId != event {
			parsed = parsed[:j+copy(parsed[j:], parsed[j+1:])]
			deleted++
		}
	}
	return parsed
}

func parseRequestArgs(vars map[string]string, varName string, writer http.ResponseWriter) (int, error) {
	id, err := strconv.Atoi(vars[varName])
	if err != nil {
		log.Println("Cannot decode request "+varName, err)
		http.Error(
			writer,
			fmt.Sprintf("Failed to decode request: %s", err),
			http.StatusBadRequest)
	}
	return id, err
}

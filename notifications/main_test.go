package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

var router *mux.Router

const publishKey = "PublishKey"

func init() {
	// Use dummy connections
	db = &dummy_db{}
	pb = &dummy_pusher_beam{}
	pc = &dummy_pusher{}

	router = mux.NewRouter()
	Init(router)
}

func TestGETNotificationWithValues(t *testing.T) {
	// Event has one entry
	req, _ := http.NewRequest("GET", "/events/55/notifications", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	expected := "[{\"title\":\"test\",\"description\":\"test\",\"regionIds\":[55,55,55],\"occurredAt\":100,\"notificationId\":55,\"eventId\":55}]"
	body := response.Body.String()
	if strings.TrimSpace(body) != expected {
		t.Errorf("Expected %s. Got %s", expected, body)
	}
}

func TestGETNotificationWithValuesSorted(t *testing.T) {
	// Event has one entry
	req, _ := http.NewRequest("GET", "/events/99/notifications", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	expected := "[{\"title\":\"test\",\"description\":\"test\",\"regionIds\":[99,99,99],\"occurredAt\":500,\"notificationId\":99,\"eventId\":99},{\"title\":\"test\",\"description\":\"test\",\"regionIds\":[99,99,99],\"occurredAt\":100,\"notificationId\":99,\"eventId\":99}]"
	body := response.Body.String()
	if strings.TrimSpace(body) != expected {
		t.Errorf("Expected %s. Got %s", expected, body)
	}
}

func TestGETNoResults(t *testing.T) {
	// Event has no entry
	req, _ := http.NewRequest("GET", "/events/44/notifications", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	expected := "[]"
	if body := response.Body.String(); strings.TrimSpace(body) != expected {
		t.Errorf("Expected an empty list. Got %s", body)
	}
}

func TestValidNotificationUpdate(t *testing.T) {
	var buf bytes.Buffer

	update := organiser_notification{
		RegionIds:   []int{99, 99, 99},
		OccurredAt:  123456,
		Title:       "title",
		Description: "description",
	}

	err := json.NewEncoder(&buf).Encode(&update)
	if err != nil {
		t.Error("Unable to encode update struct to json")
	}

	req, _ := http.NewRequest("POST", "/events/99/notifications", &buf)
	response := executeRequest(req)

	publishId := strconv.Itoa(hashString(publishKey))

	checkResponseCode(t, http.StatusOK, response.Code)
	expected := "{\"title\":\"title\",\"description\":\"description\",\"regionIds\":[99,99,99],\"occurredAt\":123456,\"notificationId\":" + publishId + ",\"eventId\":99}\n"
	body := response.Body.String()
	if body != expected {
		t.Errorf("Expected %s. Got %s", expected, body)
	}
}

func TestValidLocationUpdateWithoutDescription(t *testing.T) {
	var buf bytes.Buffer

	update := organiser_notification{
		RegionIds:  []int{99, 99, 99},
		OccurredAt: 123456,
		Title:      "title",
	}

	err := json.NewEncoder(&buf).Encode(&update)
	if err != nil {
		t.Error("Unable to encode update struct to json")
	}

	req, _ := http.NewRequest("POST", "/events/99/notifications", &buf)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestInvalidNotificationWithoutOccurredAt(t *testing.T) {
	var buf bytes.Buffer

	update := organiser_notification{
		RegionIds:   []int{99, 99, 99},
		Title:       "title",
		Description: "description",
	}

	err := json.NewEncoder(&buf).Encode(&update)
	if err != nil {
		t.Error("Unable to encode update struct to json")
	}

	req, _ := http.NewRequest("POST", "/events/99/notifications", &buf)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestInvalidNotificationWithoutRegions(t *testing.T) {
	var buf bytes.Buffer

	update := organiser_notification{
		OccurredAt:  123456,
		Title:       "title",
		Description: "description",
	}

	err := json.NewEncoder(&buf).Encode(&update)
	if err != nil {
		t.Error("Unable to encode update struct to json")
	}

	req, _ := http.NewRequest("POST", "/events/99/notifications", &buf)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestInvalidNotificationWithoutTitle(t *testing.T) {
	var buf bytes.Buffer

	update := organiser_notification{
		EventId:     99,
		RegionIds:   []int{99, 99, 99},
		OccurredAt:  123456,
		Description: "description",
	}

	err := json.NewEncoder(&buf).Encode(&update)
	if err != nil {
		t.Error("Unable to encode update struct to json")
	}

	req, _ := http.NewRequest("POST", "/events/99/notifications", &buf)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

/***************************
   FAKE DynamoDB Database
***************************/

type dummy_db struct {
	t *testing.T
}

func (dq *dummy_db) InitConn(tableName string) error {
	return nil
}

func (db *dummy_db) GetTableScan() []map[string]interface{} {
	// Make a fake table and insert a row
	tableScan := make([]map[string]interface{}, 3)
	tableScan[0] = db.makeRow(99, 100, "test")
	tableScan[1] = db.makeRow(99, 500, "test")
	tableScan[2] = db.makeRow(55, 100, "test")

	fmt.Println(tableScan)
	return tableScan
}

func (db *dummy_db) makeRow(n int, time int, s string) map[string]interface{} {
	// Use the same number, string, and bool for all values to make testing easier

	// Make the row
	row := make(map[string]interface{})
	row["eventId"] = n
	row["notificationId"] = n
	row["occurredAt"] = time

	row["title"] = s
	row["description"] = s

	// Set the boolean dealt with
	row["regionIds"] = []int{n, n, n}

	return row
}

func (db *dummy_db) SendItem(req interface{}) {
	return
}

func (db *dummy_db) GetItem(pKeyColName string, pKeyValue string) map[string]interface{} {
	return nil
}

/***************************
   FAKE Pusher queue
***************************/

type dummy_pusher struct {
	t *testing.T
}

func (pc *dummy_pusher) InitConn() {
	return
}

func (pc *dummy_pusher) SendItem(channelName string, eventName string, data []byte) {
	return
}

/***************************
   FAKE Pusher beam
***************************/

type dummy_pusher_beam struct {
	ct *testing.T
}

func (pbc *dummy_pusher_beam) InitConn() {
	return
}

func (pbc *dummy_pusher_beam) SendNotification(regionIds []string, title string, body string) (publishId string, err error) {
	return publishKey, nil
}

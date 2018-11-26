package emergency

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var router *mux.Router

func init() {
	router = mux.NewRouter()
	Init(router)
}

func TestGETLocationWithValues(t *testing.T) {
	// Create a dummy db
	db = &dummy_db{t}
	pc = &dummy_pusher{t}

	// Event has one entry
	req, _ := http.NewRequest("GET", "/live/emergency/99/0", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	expected := "[{\"uuid\":\"test\",\"eventId\":99,\"regionIds\":[99,99,99],\"occurredAt\":99,\"dealtWith\":false,\"description\":\"test\",\"position\":{\"lat\":99,\"lng\":99}}]"
	body := response.Body.String()
	if strings.TrimSpace(body) != expected {
		t.Errorf("Expected %s. Got %s", expected, body)
	}
}

func TestGETNoResults(t *testing.T) {
	// Create a dummy db with no entries
	db = &dummy_db{t}
	pc = &dummy_pusher{t}

	// Event has no entry
	req, _ := http.NewRequest("GET", "/live/emergency/1/0", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	expected := "[]"
	if body := response.Body.String(); strings.TrimSpace(body) != expected {
		t.Errorf("Expected an empty list. Got %s", body)
	}
}

func TestGETWithUpdateURL(t *testing.T) {
	req, _ := http.NewRequest("GET", "/emergency-update", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusMethodNotAllowed, response.Code)
	if body := response.Body.String(); body != "" {
		t.Errorf("Expected an empty body. Got %s", body)
	}
}

func TestValidLocationUpdate(t *testing.T) {
	var buf bytes.Buffer

	update := emergency_request{
		UUID:        "Test-UUID",
		EventId:     99,
		RegionIds:   []int{99, 99, 99},
		DealtWith:   false,
		OccurredAt:  123456,
		Description: "Help me",
	}

	err := json.NewEncoder(&buf).Encode(&update)
	if err != nil {
		t.Error("Unable to encode update struct to json")
	}

	req, _ := http.NewRequest("POST", "/emergency-update", &buf)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	expected := "{\"uuid\":\"Test-UUID\",\"eventId\":99,\"regionIds\":[99,99,99],\"occurredAt\":123456,\"dealtWith\":false,\"description\":\"Help me\",\"position\":null}\n"
	body := response.Body.String()
	if body != expected {
		t.Errorf("Expected %s. Got %s", expected, body)
	}
}

func TestValidLocationUpdateWithoutPosition(t *testing.T) {
	var buf bytes.Buffer

	update := emergency_request{
		UUID:       "Test-UUID",
		EventId:    99,
		RegionIds:  []int{99, 99, 99},
		DealtWith:  false,
		OccurredAt: 123456,
	}

	err := json.NewEncoder(&buf).Encode(&update)
	if err != nil {
		t.Error("Unable to encode update struct to json")
	}

	req, _ := http.NewRequest("POST", "/emergency-update", &buf)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	expected := "{\"uuid\":\"Test-UUID\",\"eventId\":99,\"regionIds\":[99,99,99],\"occurredAt\":123456,\"dealtWith\":false,\"description\":\"\",\"position\":null}\n"
	body := response.Body.String()
	if strings.Compare(expected, body) != 0 {
		t.Errorf("\n%s\n%s", expected, body)
	}
}

func TestValidLocationUpdateWithoutDescription(t *testing.T) {
	var buf bytes.Buffer

	update := emergency_request{
		UUID:       "Test-UUID",
		EventId:    99,
		RegionIds:  []int{99, 99, 99},
		DealtWith:  false,
		OccurredAt: 123456,
		Position: &struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		}{Lat: 1.1, Lng: 1.1},
	}

	err := json.NewEncoder(&buf).Encode(&update)
	if err != nil {
		t.Error("Unable to encode update struct to json")
	}

	req, _ := http.NewRequest("POST", "/emergency-update", &buf)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	expected := "{\"uuid\":\"Test-UUID\",\"eventId\":99,\"regionIds\":[99,99,99],\"occurredAt\":123456,\"dealtWith\":false,\"description\":\"\",\"position\":{\"lat\":1.1,\"lng\":1.1}}\n"
	if body := response.Body.String(); body != expected {
		t.Errorf("Expected %s. Got %s", expected, body)
	}
}

func TestLocationUpdateMissingTimestamp(t *testing.T) {
	var buf bytes.Buffer

	update := emergency_request{
		UUID:        "Test-UUID",
		EventId:     99,
		RegionIds:   []int{99, 99, 99},
		DealtWith:   false,
		Description: "Test-emergency",
	}

	err := json.NewEncoder(&buf).Encode(&update)
	if err != nil {
		t.Error("Unable to encode update struct to json")
	}

	req, _ := http.NewRequest("POST", "/emergency-update", &buf)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestLocationUpdateMissingUUID(t *testing.T) {
	var buf bytes.Buffer

	update := emergency_request{
		EventId:    99,
		RegionIds:  []int{99, 99, 99},
		DealtWith:  false,
		OccurredAt: int(time.Now().Unix()),
	}

	err := json.NewEncoder(&buf).Encode(&update)
	if err != nil {
		t.Error("Unable to encode update struct to json")
	}

	req, _ := http.NewRequest("POST", "/emergency-update", &buf)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestLocationUpdateMissingEventId(t *testing.T) {
	var buf bytes.Buffer

	update := emergency_request{
		UUID:        "Test-UUID",
		RegionIds:   []int{99, 99, 99},
		DealtWith:   false,
		OccurredAt:  int(time.Now().Unix()),
		Description: "",
	}

	err := json.NewEncoder(&buf).Encode(&update)
	if err != nil {
		t.Error("Unable to encode update struct to json")
	}

	req, _ := http.NewRequest("POST", "/emergency-update", &buf)
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
	tableScan := make([]map[string]interface{}, 1)
	tableScan[0] = db.makeRow(99, "test", false)

	return tableScan
}

func (db *dummy_db) makeRow(n int, s string, b bool) map[string]interface{} {
	// Use the same number, string, and bool for all values to make testing easier

	// Make the row
	row := make(map[string]interface{})
	row["eventId"] = n
	row["occurredAt"] = n

	row["uuid"] = s
	row["description"] = s

	// Set the boolean dealt with
	row["dealtWith"] = b

	// Set the boolean dealt with
	row["regionIds"] = []int{n, n, n}

	// Set the positions map
	position := make(map[string]int)
	position["lat"] = n
	position["lng"] = n
	row["position"] = position

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

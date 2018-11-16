package emergency

import (
	"bytes"
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

func (dq *dummy_db) initConn() error {
	return nil
}

func (db *dummy_db) makeDummyEntry() {

}

func (db *dummy_db) getTableScan() (*dynamodb.ScanOutput, error) {
	// Make a fake table and insert a row
	tableScan := new(dynamodb.ScanOutput)
	row := db.makeDynamoDbRow("99", "test", false)
	tableScan.Items = append(tableScan.Items, row)

	return tableScan, nil
}

func (db *dummy_db) makeDynamoDbRow(n string, s string, b bool) map[string]*dynamodb.AttributeValue {
	// Use the same number, string, and bool for all values to make testing easier

	// Make the row
	row := make(map[string]*dynamodb.AttributeValue)

	// Set the uuid & occurred at (numerical values)
	eventId := dynamodb.AttributeValue{}
	eventId.N = &n
	row["eventId"] = &eventId
	row["occurredAt"] = &eventId

	// Set the uuid
	uuid := dynamodb.AttributeValue{}
	uuid.S = &s
	row["uuid"] = &uuid
	row["description"] = &uuid

	// Set the boolean dealt with
	dealtWith := dynamodb.AttributeValue{}
	dealtWith.BOOL = &b
	row["dealtWith"] = &dealtWith

	// Set the boolean dealt with
	regionIds := dynamodb.AttributeValue{}
	regionIds.L = []*dynamodb.AttributeValue{&eventId, &eventId, &eventId}
	row["regionIds"] = &regionIds

	// Set the positions map
	row["position"] = &dynamodb.AttributeValue{}
	row["position"].M = make(map[string]*dynamodb.AttributeValue)
	(row["position"].M)["lat"] = &eventId
	(row["position"].M)["lng"] = &eventId

	return row
}

func (db *dummy_db) sendItem(req emergency_request) (err error) {
	return nil
}

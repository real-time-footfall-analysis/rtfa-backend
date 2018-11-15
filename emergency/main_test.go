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
	expected := "[{\"uuid\":\"test\",\"eventId\":99,\"regionIds\":[99,99,99],\"occurredAt\":99,\"dealtWith\":false,\"description\":\"test\"}]"
	if body := response.Body.String(); strings.TrimSpace(body) != expected {
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
		t.Errorf("Expected an empty body. Got %s", body)
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
		UUID:       "Test-UUID",
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

	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "" {
		t.Errorf("Expected an empty body. Got %s", body)
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

	return row
}

func (db *dummy_db) sendItem(req emergency_request) (err error) {
	return nil
}

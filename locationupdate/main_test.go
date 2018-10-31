package locationupdate

import (
	"bytes"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

var router *mux.Router

func init() {
	router = mux.NewRouter()
	Init(router)
}

func TestGETLocationUpdate(t *testing.T) {
	req, _ := http.NewRequest("GET", "/update", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusMethodNotAllowed, response.Code)
	if body := response.Body.String(); body != "" {
		t.Errorf("Expected an empty body. Got %s", body)
	}
}

func TestEmptyLocationUpdate(t *testing.T) {
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	queue = &dummy_queue{update: update{}, t: t}

	var buf bytes.Buffer
	req, _ := http.NewRequest("POST", "/update", &buf)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
	if body := response.Body.String(); strings.TrimSpace(body) != "Failed to decode update: EOF" {
		t.Errorf("Expected \"Failed to decode update: EOF\". Got \"%s\"", body)
	}

	r, _ := regexp.Compile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} Cannot decode update: EOF\n$`)
	if !r.MatchString(logBuf.String()) {
		t.Error("Expected Log output for failed decoding of empty body")
	}
}

func TestUUIDLengthLocationUpdate(t *testing.T) {
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	queue = &dummy_queue{update: update{}, t: t}

	var buf bytes.Buffer
	buf.WriteString(`{"uuid":"UUID", "eventId":0,"regionId":1,"entering":true,"occurredAt":1540945705}`)
	req, _ := http.NewRequest("POST", "/update", &buf)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
	if body := response.Body.String(); strings.TrimSpace(body) != "UUID less than "+strconv.Itoa(UUID_MIN_LENGTH)+" characters" {
		t.Errorf("Expected \"UUID less than "+strconv.Itoa(UUID_MIN_LENGTH)+" characters\". Got \"%s\"", body)
	}

	r, _ := regexp.Compile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} UUID less than ` + strconv.Itoa(UUID_MIN_LENGTH) + ` characters in update\n$`)
	if !r.MatchString(logBuf.String()) {
		t.Error("Expected Log output for UUID less than " + strconv.Itoa(UUID_MIN_LENGTH) + " characters")
	}
}

func TestIncompleteLocationUpdate1(t *testing.T) {
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	queue = &dummy_queue{update: update{}, t: t}

	var buf bytes.Buffer
	buf.WriteString(`{"eventId":0,"regionId":1,"entering":true,"occurredAt":1540945705}`)
	req, _ := http.NewRequest("POST", "/update", &buf)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
	expected := "UUID not present in update"
	if body := response.Body.String(); strings.TrimSpace(body) != expected {
		t.Errorf("Expected \"%s\". Got \"%s\"", expected, body)
	}

	r, _ := regexp.Compile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} ` + expected + `\n$`)
	if !r.MatchString(logBuf.String()) {
		t.Errorf("Expected Log output for %s", expected)
	}
}

func TestIncompleteLocationUpdate2(t *testing.T) {
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	queue = &dummy_queue{update: update{}, t: t}

	var buf bytes.Buffer
	buf.WriteString(`{"uuid":"Test-UUID","regionId":1,"entering":true,"occurredAt":1540945705}`)
	req, _ := http.NewRequest("POST", "/update", &buf)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
	expected := "EventID not present in update"
	if body := response.Body.String(); strings.TrimSpace(body) != expected {
		t.Errorf("Expected \"%s\". Got \"%s\"", expected, body)
	}

	r, _ := regexp.Compile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} ` + expected + `\n$`)
	if !r.MatchString(logBuf.String()) {
		t.Errorf("Expected Log output for %s", expected)
	}
}

func TestIncompleteLocationUpdate3(t *testing.T) {
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	queue = &dummy_queue{update: update{}, t: t}

	var buf bytes.Buffer
	buf.WriteString(`{"uuid":"Test-UUID", "eventId":0,"entering":true,"occurredAt":1540945705}`)
	req, _ := http.NewRequest("POST", "/update", &buf)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
	expected := "RegionID not present in update"
	if body := response.Body.String(); strings.TrimSpace(body) != expected {
		t.Errorf("Expected \"%s\". Got \"%s\"", expected, body)
	}

	r, _ := regexp.Compile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} ` + expected + `\n$`)
	if !r.MatchString(logBuf.String()) {
		t.Errorf("Expected Log output for %s", expected)
	}
}

func TestIncompleteLocationUpdate4(t *testing.T) {
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	queue = &dummy_queue{update: update{}, t: t}

	var buf bytes.Buffer
	buf.WriteString(`{"uuid":"Test-UUID","eventId":0,"regionId":1,"occurredAt":1540945705}`)
	req, _ := http.NewRequest("POST", "/update", &buf)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
	expected := "Entering not present in update"
	if body := response.Body.String(); strings.TrimSpace(body) != expected {
		t.Errorf("Expected \"%s\". Got \"%s\"", expected, body)
	}

	r, _ := regexp.Compile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} ` + expected + `\n$`)
	if !r.MatchString(logBuf.String()) {
		t.Errorf("Expected Log output for %s", expected)
	}
}

func TestIncompleteLocationUpdate5(t *testing.T) {
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	queue = &dummy_queue{update: update{}, t: t}

	var buf bytes.Buffer
	buf.WriteString(`{"uuid":"Test-UUID","eventId":0,"regionId":1,"entering":true}`)
	req, _ := http.NewRequest("POST", "/update", &buf)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
	expected := "OccurredAt not present in update"
	if body := response.Body.String(); strings.TrimSpace(body) != expected {
		t.Errorf("Expected \"%s\". Got \"%s\"", expected, body)
	}

	r, _ := regexp.Compile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} ` + expected + `\n$`)
	if !r.MatchString(logBuf.String()) {
		t.Errorf("Expected Log output for %s", expected)
	}
}

func TestLocationUpdate(t *testing.T) {
	var buf bytes.Buffer

	uuid := "Test-UUID"
	eventId := 0
	regionID := 1
	entering := true
	occurredAt := int(time.Now().Unix())
	update := update{
		UUID:       &uuid,
		EventID:    &eventId,
		RegionID:   &regionID,
		Entering:   &entering,
		OccurredAt: &occurredAt,
	}

	queue = &dummy_queue{update: update, t: t}

	err := json.NewEncoder(&buf).Encode(&update)
	if err != nil {
		t.Error("Unable to encode update struct to json")
	}

	req, _ := http.NewRequest("POST", "/update", &buf)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "" {
		t.Errorf("Expected an empty body. Got %s", body)
	}
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

type dummy_queue struct {
	update update
	t      *testing.T
}

// initConn opens the connection to the location event kinesis queue
func (dq *dummy_queue) initConn() error {
	return nil
}

// Pre: the event object is valid
func (dq *dummy_queue) addLocationUpdate(event *update) error {
	match := true
	if *event.UUID != *dq.update.UUID {
		match = false
	}
	if *event.EventID != *dq.update.EventID {
		match = false
	}
	if *event.RegionID != *dq.update.RegionID {
		match = false
	}
	if *event.Entering != *dq.update.Entering {
		match = false
	}
	if math.Abs(float64(*event.OccurredAt-*dq.update.OccurredAt)) > 2 {
		match = false
	}

	if !match {
		dq.t.Error("incorrect update fired to queue")
	}
	return nil
}

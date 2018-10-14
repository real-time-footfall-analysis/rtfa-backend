package locationupdate

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
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

func TestLocationUpdate(t *testing.T) {
	var buf bytes.Buffer
	update := BluetoothUpdate{
		UUID:        "Test-UUID",
		EventID:     "Test-Event",
		BeaconMajor: 7357,
		BeaconMinor: 1234,
		Entering:    true,
	}

	err := json.NewEncoder(&buf).Encode(&update)
	if err != nil {
		t.Error("Unable to encode BluetoothUpdate struct to json")
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

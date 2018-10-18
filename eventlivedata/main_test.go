package eventlivedata

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var router *mux.Router

func init() {
	router = mux.NewRouter()
	Init(router)

}

func TestGETLocationUpdate(t *testing.T) {

	retMap := make(map[string]int, 0)
	retMap["one"] = 1
	db = &dummy_db{retMap, 1, t}

	req, _ := http.NewRequest("GET", "/live/heatmap/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); strings.TrimSpace(body) != "{\"one\":1}" {
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

type dummy_db struct {
	ret map[string]int
	id  int
	t   *testing.T
}

// initConn opens the connection to the location event kinesis queue
func (dq *dummy_db) initConn() error {
	return nil
}

// Pre: the event object is valid
func (db *dummy_db) getLiveHeatMap(event int) (map[string]int, error) {
	if db.id != event {
		db.t.Errorf("wrong event id")
	}
	return db.ret, nil
}

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
	db = &dummy_db{t}

	req, _ := http.NewRequest("GET", "/live/heatmap/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	expected := "{\"1\":1}"
	if body := response.Body.String(); strings.TrimSpace(body) != expected {
		t.Errorf("Expected %s. Got %s", expected, body)
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
	t *testing.T
}

func (dq *dummy_db) InitConn(tableName string) error {
	return nil
}

func (db *dummy_db) GetTableScan() []map[string]interface{} {
	row := make(map[string]interface{})
	row["eventId"] = 1
	row["regionId"] = 1
	rows := make([]map[string]interface{}, 1)
	rows[0] = row
	return rows
}

func (db *dummy_db) SendItem(req interface{}) {
	return
}

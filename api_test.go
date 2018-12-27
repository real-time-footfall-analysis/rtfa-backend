package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a = App{}

func init() {
	initialize(&a)
}

func TestServerHealth(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/health", nil)
	response := executeRequest(a, req)

	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "" {
		t.Errorf("Expected an empty body. Got %s", body)
	}
}

type HelloWorldResponse struct {
	Message string `json:"message"`
}

func TestServerDefaultPath(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	response := executeRequest(a, req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var helloWorldResponse HelloWorldResponse

	err := json.NewDecoder(response.Body).Decode(&helloWorldResponse)
	if err != nil {
		t.Errorf("Expected json, got decode error")
	}
	if helloWorldResponse.Message != "Hello, World!" {
		t.Errorf("Expected Hello World. Got %s",
			helloWorldResponse.Message)
	}
}

func executeRequest(a App, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// create test for VersionHandler function

func TestVersionHandler(t *testing.T) {
	// create request
	req, err := http.NewRequest("GET", "/version", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(VersionHandler)

	// call handler
	handler.ServeHTTP(rr, req)

	// check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// check response
	expected := "App Version: 0.0.1"

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

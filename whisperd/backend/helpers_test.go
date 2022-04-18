package backend

import (
	"encoding/json"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"
)

func assertStatus(t *testing.T, rec *httptest.ResponseRecorder, expected int) {
	if rec.Result().StatusCode != expected {
		t.Fatalf("Expected status code %v but got %v", expected, rec.Result().StatusCode)
	}
}

func assertContentType(t *testing.T, rec *httptest.ResponseRecorder, expected string) {
	ct := rec.Result().Header.Get("Content-type")
	if ct != expected {
		t.Fatalf("Expected content type to be '%v' but got '%v'", expected, ct)
	}
}

func assertValidResponse(t *testing.T, msgMatch regexp.Regexp) func(rec *httptest.ResponseRecorder) {
	return func(rec *httptest.ResponseRecorder) {
		var s shoutResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &s); err != nil {
			t.Fatal(err)
		}
		if s.Shout.ID == 0 {
			t.Fatalf("Expected response to have non-zero id but got %d", s.Shout.ID)
		}
		if !msgMatch.Match([]byte(s.Shout.Message)) {
			t.Fatalf("Expected response message to match '%s' but was '%s'", msgMatch.String(), s.Shout.Message)
		}
		if (s.Shout.Timestamp == time.Time{}) {
			t.Fatalf("Expected response timestamp to not be empty but was '%s'", s.Shout.Timestamp)
		}
		if s.Shout.Timestamp.After(time.Now()) {
			t.Fatalf("Expected response timestamp to not be in the future but was '%s'", s.Shout.Timestamp)
		}
	}
}

func assertErrorResponse(t *testing.T, r regexp.Regexp) func(rec *httptest.ResponseRecorder) {
	return func(rec *httptest.ResponseRecorder) {
		var e errorResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &e); err != nil {
			t.Fatal(err)
		}
		if !r.Match([]byte(e.Error)) {
			t.Fatalf("Expected error message to match '%s' but got '%s'", r.String(), e.Error)
		}
	}
}

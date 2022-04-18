package backend

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"regexp"
	"time"
)

func assertStatus(rec *httptest.ResponseRecorder, expected int) error {
	if rec.Result().StatusCode != expected {
		return fmt.Errorf("Expected status code %v but got %v", expected, rec.Result().StatusCode)
	}
	return nil
}

func assertContentType(rec *httptest.ResponseRecorder, expected string) error {
	ct := rec.Result().Header.Get("Content-type")
	if ct != expected {
		return fmt.Errorf("Expected content type to be '%v' but got '%v'", expected, ct)
	}
	return nil
}

func assertValidResponse(msgMatch regexp.Regexp) func(rec *httptest.ResponseRecorder) error {
	return func(rec *httptest.ResponseRecorder) error {
		var s shoutResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &s); err != nil {
			return err
		}
		if s.Shout.ID == 0 {
			return fmt.Errorf("Expected response to have non-zero id but got %d", s.Shout.ID)
		}
		if !msgMatch.Match([]byte(s.Shout.Message)) {
			return fmt.Errorf("Expected response message to match '%s' but was '%s'", msgMatch.String(), s.Shout.Message)
		}
		if (s.Shout.Timestamp == time.Time{}) {
			return fmt.Errorf("Expected response timestamp to not be empty but was '%s'", s.Shout.Timestamp)
		}
		if s.Shout.Timestamp.After(time.Now()) {
			return fmt.Errorf("Expected response timestamp to not be in the future but was '%s'", s.Shout.Timestamp)
		}
		return nil
	}
}

func assertErrorResponse(r regexp.Regexp) func(rec *httptest.ResponseRecorder) error {
	return func(rec *httptest.ResponseRecorder) error {
		var e errorResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &e); err != nil {
			return err
		}
		if !r.Match([]byte(e.Error)) {
			return fmt.Errorf("Expected error message to match '%s' but got '%s'", r.String(), e.Error)
		}
		return nil
	}
}

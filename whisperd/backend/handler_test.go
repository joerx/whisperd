package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"
	"time"

	"whisperd.io/whisperd/whisperd"
	"whisperd.io/whisperd/whisperd/db"
)

const contentType = "application/json"

type shoutStoreMock struct {
	shouts []whisperd.Shout
}

// GetAll implements db.ShoutStore.GetAll for tests
func (m *shoutStoreMock) GetAll(ctx context.Context) ([]whisperd.Shout, error) {
	return nil, nil
}

// Get implements db.ShoutStore.Get for tests
func (m *shoutStoreMock) Get(ctx context.Context, id string) (whisperd.Shout, error) {
	tmp, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
		return whisperd.Shout{}, fmt.Errorf("failed to parse id '%s' as integer", id)
	}

	i := tmp - 1

	if i < 0 {
		return whisperd.Shout{}, fmt.Errorf("invalid id '%s'", id)
	}
	if int(i) >= len(m.shouts) {
		return whisperd.Shout{}, db.ErrorNotFound
	}
	return m.shouts[i], nil
}

// Put implements db.ShoutStore.Put for tests
func (m *shoutStoreMock) Insert(ctx context.Context, s whisperd.Shout) (whisperd.Shout, error) {
	if s.Message == "" {
		return s, db.ErrorInvalidRecord
	}
	return whisperd.Shout{ID: 1, Message: s.Message, Timestamp: time.Now()}, nil
}

// Delete implements db.ShoutStore.Delete for tests
func (m *shoutStoreMock) Delete(ctx context.Context, s whisperd.Shout) (whisperd.Shout, error) {
	return s, nil
}

func testHandler() http.Handler {
	m := &shoutStoreMock{
		shouts: []whisperd.Shout{
			{ID: 1, Message: "Hello World", Timestamp: time.Now()},
		},
	}

	h := newHandler(m)
	rt := newRouter(h)

	return rt
}

func TestGetShout(t *testing.T) {
	type testCase struct {
		url    string
		sc     int
		assert func(*httptest.ResponseRecorder)
	}

	testCases := []testCase{
		{
			url:    "/api/shout/1",
			sc:     http.StatusOK,
			assert: assertValidResponse(t, *regexp.MustCompile("Hello World")),
		},
		{
			url:    "/api/shout/20",
			sc:     http.StatusNotFound,
			assert: assertErrorResponse(t, *regexp.MustCompile("not found")),
		},
	}

	rt := testHandler()

	for _, tc := range testCases {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, tc.url, nil)

		rt.ServeHTTP(rec, req)

		assertStatus(t, rec, tc.sc)
		assertContentType(t, rec, contentType)
		tc.assert(rec)
	}
}

func TestPostShout(t *testing.T) {
	type testCase struct {
		input  putShoutRequest
		status int
		assert func(*httptest.ResponseRecorder)
	}

	tests := []testCase{
		{
			input:  putShoutRequest{Shout: whisperd.Shout{Message: "sshhh!"}},
			status: http.StatusOK,
			assert: assertValidResponse(t, *regexp.MustCompile("sshhh!")),
		},
		{
			input:  putShoutRequest{Shout: whisperd.Shout{}},
			status: http.StatusBadRequest,
			assert: assertErrorResponse(t, *regexp.MustCompile("invalid record")),
		},
	}

	rt := testHandler()

	for _, tc := range tests {
		data, err := json.Marshal(tc.input)
		if err != nil {
			t.Fatal(err)
		}

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/shouts", bytes.NewReader(data))

		rt.ServeHTTP(rec, req)

		assertStatus(t, rec, tc.status)
		assertContentType(t, rec, contentType)
		tc.assert(rec)
	}
}

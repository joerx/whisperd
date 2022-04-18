package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
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
	return m.shouts, nil
}

// Get implements db.ShoutStore.Get for tests
func (m *shoutStoreMock) Get(ctx context.Context, id int64) (whisperd.Shout, error) {
	if id < 1 {
		return whisperd.Shout{}, fmt.Errorf("invalid id '%d'", id)
	}
	if int(id-1) > len(m.shouts) {
		return whisperd.Shout{}, db.ErrorNotFound
	}
	return m.shouts[id-1], nil
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

func testHandler(list []whisperd.Shout) http.Handler {
	m := &shoutStoreMock{
		shouts: list,
	}

	h := newHandler(m)
	rt := newRouter(h)

	return rt
}

func TestGetShout(t *testing.T) {
	shouts := []whisperd.Shout{
		{ID: 1, Message: "Moo!", Timestamp: time.Now()},
	}
	rt := testHandler(shouts)

	type testCase struct {
		url    string
		sc     int
		assert func(*httptest.ResponseRecorder) error
	}

	testCases := []testCase{
		{
			url:    "/api/shout/1",
			sc:     http.StatusOK,
			assert: assertValidResponse(*regexp.MustCompile(shouts[0].Message)),
		},
		{
			url:    "/api/shout/20",
			sc:     http.StatusNotFound,
			assert: assertErrorResponse(*regexp.MustCompile("not found")),
		},
	}

	for _, tc := range testCases {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, tc.url, nil)

		rt.ServeHTTP(rec, req)

		if err := assertStatus(rec, tc.sc); err != nil {
			t.Fatal(err)
		}
		if err := assertContentType(rec, contentType); err != nil {
			t.Fatal(err)
		}
		if err := tc.assert(rec); err != nil {
			t.Fatal(err)
		}
	}
}

func TestListShouts(t *testing.T) {
	type test struct {
		shouts []whisperd.Shout
		status int
	}

	tests := []test{
		{
			status: http.StatusOK,
			shouts: []whisperd.Shout{
				{ID: 1, Message: "Halloballo", Timestamp: time.Now()},
				{ID: 2, Message: "Elon is an ugly cow!", Timestamp: time.Now()},
			},
		},
		{
			status: http.StatusOK,
			shouts: []whisperd.Shout{},
		},
	}

	for _, tc := range tests {
		rt := testHandler(tc.shouts)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/shouts", nil)

		rt.ServeHTTP(rec, req)

		var rl shoutListResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &rl); err != nil {
			t.Fatal(err)
		}

		if err := assertStatus(rec, tc.status); err != nil {
			t.Fatal(err)
		}
		if err := assertContentType(rec, contentType); err != nil {
			t.Fatal(err)
		}
		if len(rl.Shouts) != len(tc.shouts) {
			t.Fatalf("expected %d shouts to be returned but got %d", len(tc.shouts), len(rl.Shouts))
		}
	}
}

func TestPostShout(t *testing.T) {
	type testCase struct {
		input  putShoutRequest
		status int
		assert func(*httptest.ResponseRecorder) error
	}

	tests := []testCase{
		{
			input:  putShoutRequest{Shout: whisperd.Shout{Message: "sshhh!"}},
			status: http.StatusOK,
			assert: assertValidResponse(*regexp.MustCompile("sshhh!")),
		},
		{
			input:  putShoutRequest{Shout: whisperd.Shout{}},
			status: http.StatusBadRequest,
			assert: assertErrorResponse(*regexp.MustCompile("invalid record")),
		},
	}

	rt := testHandler([]whisperd.Shout{})

	for _, tc := range tests {
		data, err := json.Marshal(tc.input)
		if err != nil {
			t.Fatal(err)
		}

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/shouts", bytes.NewReader(data))

		rt.ServeHTTP(rec, req)

		if err := assertStatus(rec, tc.status); err != nil {
			t.Fatal(err)
		}
		if err := assertContentType(rec, contentType); err != nil {
			t.Fatal(err)
		}
		if err := tc.assert(rec); err != nil {
			t.Fatal(err)
		}
	}
}

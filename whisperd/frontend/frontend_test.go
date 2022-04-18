package frontend

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFrontend(t *testing.T) {
	type test struct {
		path string
		code int
		ct   string
	}

	tests := []test{
		{
			path: "/index.html",
			code: http.StatusOK,
			ct:   "text/html; charset=utf-8",
		},
		{
			path: "/",
			code: http.StatusOK,
			ct:   "text/html; charset=utf-8",
		},
		{
			path: "/foo",
			code: http.StatusNotFound,
			ct:   "text/plain",
		},
		{
			path: "/css",
			code: http.StatusNotFound,
			ct:   "text/plain",
		},
		{
			path: "/css/index.css",
			code: http.StatusOK,
			ct:   "text/plain; charset=utf-8",
		},
	}

	h, err := Handler()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range tests {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		sc := rec.Result().StatusCode
		ct := rec.Result().Header.Get("Content-type")

		if sc != tc.code {
			t.Fatalf("Expected status code for %s to be %d but got %d", tc.path, tc.code, sc)
		}
		if ct != tc.ct {
			t.Fatalf("Expected content type for %s to be '%s' but was '%s'", tc.path, tc.ct, ct)
		}
	}
}

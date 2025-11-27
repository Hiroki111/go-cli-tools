package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupAPI(t *testing.T) (string, func()) {
	t.Helper()

	server := httptest.NewServer(newMux(""))

	return server.URL, func() {
		server.Close()
	}
}

func TestGet(t *testing.T) {
	testCases := []struct {
		name            string
		path            string
		expectedCode    int
		expectedItems   int
		expectedContent string
	}{
		{name: "GetRoot", path: "/", expectedCode: http.StatusOK, expectedContent: "There's an API here"},
		{name: "NotFound", path: "/todo/500", expectedCode: http.StatusNotFound},
	}

	url, cleanup := setupAPI(t)
	defer cleanup()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var (
				body []byte
				err  error
			)

			r, err := http.Get(url + testCase.path)
			if err != nil {
				t.Error(err)
			}
			defer r.Body.Close()

			if r.StatusCode != testCase.expectedCode {
				t.Fatalf("Expected %q, got %q.", http.StatusText(testCase.expectedCode), http.StatusText(r.StatusCode))
			}

			switch {
			case strings.Contains(r.Header.Get("Content-Type"), "text/plain"):
				if body, err = io.ReadAll(r.Body); err != nil {
					t.Error(err)
				}
				if !strings.Contains(string(body), testCase.expectedContent) {
					t.Errorf("Expected %q, got %q.", testCase.expectedContent, string(body))
				}
			default:
				t.Fatalf("Unsupported Content-Type: %q", r.Header.Get("Content-Type"))
			}
		})
	}
}

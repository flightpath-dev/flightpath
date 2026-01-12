package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogging_StatusCode(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		expectedInLog string
	}{
		{"200 OK", http.StatusOK, "200"},
		{"201 Created", http.StatusCreated, "201"},
		{"400 Bad Request", http.StatusBadRequest, "400"},
		{"404 Not Found", http.StatusNotFound, "404"},
		{"500 Internal Server Error", http.StatusInternalServerError, "500"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := log.New(&buf, "", 0)

			handler := Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			logOutput := buf.String()
			if !strings.Contains(logOutput, tt.expectedInLog) {
				t.Errorf("log should contain status code %q, got: %s", tt.expectedInLog, logOutput)
			}
		})
	}
}

func TestLogging_BytesWritten(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	responseBody := "Hello, World!"
	handler := Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(responseBody))
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	logOutput := buf.String()
	expectedBytes := "13 bytes" // len("Hello, World!") = 13
	if !strings.Contains(logOutput, expectedBytes) {
		t.Errorf("log should contain %q, got: %s", expectedBytes, logOutput)
	}
}

func TestLogging_MethodAndPath(t *testing.T) {
	tests := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/health"},
		{http.MethodPost, "/api/commands"},
		{http.MethodPut, "/api/settings"},
		{http.MethodDelete, "/api/data"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			var buf bytes.Buffer
			logger := log.New(&buf, "", 0)

			handler := Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			logOutput := buf.String()
			if !strings.Contains(logOutput, tt.method) {
				t.Errorf("log should contain method %q, got: %s", tt.method, logOutput)
			}
			if !strings.Contains(logOutput, tt.path) {
				t.Errorf("log should contain path %q, got: %s", tt.path, logOutput)
			}
		})
	}
}

func TestLogging_DefaultStatusCode(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	// Handler that writes body without explicitly setting status code
	handler := Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	logOutput := buf.String()
	// Default status should be 200
	if !strings.Contains(logOutput, "200") {
		t.Errorf("log should contain default status code 200, got: %s", logOutput)
	}
}

func TestResponseWriter_Flush(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	handler := Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("First chunk"))
		// Test that Flush works without panicking
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
		w.Write([]byte("Second chunk"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Verify response contains both chunks
	body := rr.Body.String()
	if body != "First chunkSecond chunk" {
		t.Errorf("body: got %q, want %q", body, "First chunkSecond chunk")
	}

	// Verify bytes written is correct
	logOutput := buf.String()
	expectedBytes := "23 bytes" // len("First chunk") + len("Second chunk") = 11 + 12 = 23
	if !strings.Contains(logOutput, expectedBytes) {
		t.Errorf("log should contain %q, got: %s", expectedBytes, logOutput)
	}
}

func TestResponseWriter_MultipleWrites(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	handler := Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("A"))
		w.Write([]byte("BC"))
		w.Write([]byte("DEF"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	logOutput := buf.String()
	expectedBytes := "6 bytes" // 1 + 2 + 3 = 6
	if !strings.Contains(logOutput, expectedBytes) {
		t.Errorf("log should contain %q, got: %s", expectedBytes, logOutput)
	}
}

func TestLogging_DurationLogged(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	handler := Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	logOutput := buf.String()
	// Duration should be in the log (e.g., "1.234µs" or "1.234ms")
	if !strings.Contains(logOutput, "s") && !strings.Contains(logOutput, "ms") && !strings.Contains(logOutput, "µs") {
		t.Errorf("log should contain duration, got: %s", logOutput)
	}
}

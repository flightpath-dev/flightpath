package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/flightpath-dev/flightpath/internal/logger"
)

func TestRecovery_NoPanic(t *testing.T) {
	var buf bytes.Buffer
	slogHandler := slog.NewTextHandler(&buf, nil)
	logger := &logger.Logger{Logger: slog.New(slogHandler)}

	handler := Recovery(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Request should be processed normally
	if rr.Code != http.StatusOK {
		t.Errorf("Status code: got %d, want %d", rr.Code, http.StatusOK)
	}
	if rr.Body.String() != "OK" {
		t.Errorf("Body: got %q, want %q", rr.Body.String(), "OK")
	}

	// No panic log should be present
	if strings.Contains(buf.String(), "PANIC") {
		t.Errorf("log should not contain PANIC for normal request, got: %s", buf.String())
	}
}

func TestRecovery_WithPanic(t *testing.T) {
	var buf bytes.Buffer
	slogHandler := slog.NewTextHandler(&buf, nil)
	logger := &logger.Logger{Logger: slog.New(slogHandler)}

	handler := Recovery(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong!")
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	// This should NOT panic - recovery middleware should catch it
	handler.ServeHTTP(rr, req)

	// Should return 500 Internal Server Error
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Status code: got %d, want %d", rr.Code, http.StatusInternalServerError)
	}

	// Response body should indicate error
	if !strings.Contains(rr.Body.String(), "Internal server error") {
		t.Errorf("Body should contain 'Internal server error', got: %q", rr.Body.String())
	}

	// Log should contain panic details
	logOutput := buf.String()
	if !strings.Contains(logOutput, "PANIC") {
		t.Errorf("log should contain PANIC, got: %s", logOutput)
	}
	if !strings.Contains(logOutput, "something went wrong!") {
		t.Errorf("log should contain panic message, got: %s", logOutput)
	}
}

func TestRecovery_WithPanicError(t *testing.T) {
	var buf bytes.Buffer
	slogHandler := slog.NewTextHandler(&buf, nil)
	logger := &logger.Logger{Logger: slog.New(slogHandler)}

	handler := Recovery(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var slice []int
		// This will panic with index out of range
		_ = slice[10]
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	// This should NOT panic - recovery middleware should catch it
	handler.ServeHTTP(rr, req)

	// Should return 500 Internal Server Error
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Status code: got %d, want %d", rr.Code, http.StatusInternalServerError)
	}

	// Log should contain stack trace
	logOutput := buf.String()
	if !strings.Contains(logOutput, "PANIC") {
		t.Errorf("log should contain PANIC, got: %s", logOutput)
	}
}

func TestRecovery_WithNilPanic(t *testing.T) {
	var buf bytes.Buffer
	slogHandler := slog.NewTextHandler(&buf, nil)
	logger := &logger.Logger{Logger: slog.New(slogHandler)}

	handler := Recovery(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ptr *string
		// This will panic with nil pointer dereference
		_ = *ptr
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	// This should NOT panic - recovery middleware should catch it
	handler.ServeHTTP(rr, req)

	// Should return 500 Internal Server Error
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Status code: got %d, want %d", rr.Code, http.StatusInternalServerError)
	}
}

func TestRecovery_StackTraceLogged(t *testing.T) {
	var buf bytes.Buffer
	slogHandler := slog.NewTextHandler(&buf, nil)
	logger := &logger.Logger{Logger: slog.New(slogHandler)}

	handler := Recovery(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	logOutput := buf.String()

	// Stack trace should contain goroutine info
	if !strings.Contains(logOutput, "goroutine") {
		t.Errorf("log should contain stack trace with goroutine info, got: %s", logOutput)
	}
}

func TestRecovery_ChainedMiddleware(t *testing.T) {
	var buf bytes.Buffer
	slogHandler := slog.NewTextHandler(&buf, nil)
	logger := &logger.Logger{Logger: slog.New(slogHandler)}

	// Chain recovery with another middleware
	handler := Recovery(logger)(
		Logging(logger)(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("chained panic")
			}),
		),
	)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	// This should NOT panic
	handler.ServeHTTP(rr, req)

	// Should return 500 Internal Server Error
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Status code: got %d, want %d", rr.Code, http.StatusInternalServerError)
	}
}

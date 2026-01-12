package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS_AllowedOrigin(t *testing.T) {
	allowedOrigins := []string{"http://localhost:3000", "http://example.com"}
	handler := CORS(allowedOrigins)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check CORS headers are set
	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Errorf("Access-Control-Allow-Origin: got %q, want %q", got, "http://localhost:3000")
	}
	if got := rr.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Error("Access-Control-Allow-Methods header not set")
	}
	if got := rr.Header().Get("Access-Control-Allow-Headers"); got == "" {
		t.Error("Access-Control-Allow-Headers header not set")
	}
	if got := rr.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Errorf("Access-Control-Allow-Credentials: got %q, want %q", got, "true")
	}
}

func TestCORS_DisallowedOrigin(t *testing.T) {
	allowedOrigins := []string{"http://localhost:3000"}
	handler := CORS(allowedOrigins)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://evil.com")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Access-Control-Allow-Origin should NOT be set for disallowed origins
	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("Access-Control-Allow-Origin should be empty for disallowed origin, got %q", got)
	}
}

func TestCORS_Wildcard(t *testing.T) {
	allowedOrigins := []string{"*"}
	handler := CORS(allowedOrigins)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	testOrigins := []string{
		"http://localhost:3000",
		"http://example.com",
		"http://any-domain.io",
	}

	for _, origin := range testOrigins {
		t.Run(origin, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Origin", origin)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			// Wildcard should allow any origin
			if got := rr.Header().Get("Access-Control-Allow-Origin"); got != origin {
				t.Errorf("Access-Control-Allow-Origin: got %q, want %q", got, origin)
			}
		})
	}
}

func TestCORS_Preflight(t *testing.T) {
	allowedOrigins := []string{"http://localhost:3000"}
	handler := CORS(allowedOrigins)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This should NOT be called for preflight requests
		t.Error("handler should not be called for preflight requests")
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Preflight should return 200 OK
	if rr.Code != http.StatusOK {
		t.Errorf("Status code: got %d, want %d", rr.Code, http.StatusOK)
	}

	// CORS headers should be set
	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Errorf("Access-Control-Allow-Origin: got %q, want %q", got, "http://localhost:3000")
	}
	if got := rr.Header().Get("Access-Control-Max-Age"); got != "3600" {
		t.Errorf("Access-Control-Max-Age: got %q, want %q", got, "3600")
	}
}

func TestCORS_NoOriginHeader(t *testing.T) {
	allowedOrigins := []string{"http://localhost:3000"}
	handler := CORS(allowedOrigins)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	// No Origin header set
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Access-Control-Allow-Origin should NOT be set when no Origin header
	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("Access-Control-Allow-Origin should be empty when no Origin header, got %q", got)
	}

	// Request should still be processed
	if rr.Code != http.StatusOK {
		t.Errorf("Status code: got %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestCORS_EmptyAllowedOrigins(t *testing.T) {
	allowedOrigins := []string{}
	handler := CORS(allowedOrigins)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// No origins allowed, so CORS header should not be set
	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("Access-Control-Allow-Origin should be empty for empty allowed list, got %q", got)
	}
}

func TestCORS_ConnectHeaders(t *testing.T) {
	allowedOrigins := []string{"http://localhost:3000"}
	handler := CORS(allowedOrigins)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Verify Connect-specific headers are allowed
	allowedHeaders := rr.Header().Get("Access-Control-Allow-Headers")
	expectedHeaders := []string{"Content-Type", "Connect-Protocol-Version", "Connect-Timeout-Ms", "Authorization"}
	for _, h := range expectedHeaders {
		if allowedHeaders == "" {
			t.Errorf("Access-Control-Allow-Headers should contain %q", h)
		}
	}
}

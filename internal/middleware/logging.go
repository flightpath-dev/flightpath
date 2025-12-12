package middleware

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

// ----------------------------------------------------------------------------
// Creates a request logging middleware that logs all HTTP requests, specifically
// their method, path, status code, duration, bytes written, and cancellation status.
// The log is written after the handler completes (including canceled requests).
// ----------------------------------------------------------------------------
func Logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now() // Track request start time for duration calculation

			// Wrap response writer to capture status code and bytes written
			wrapped := newResponseWriter(w)

			// Process request - this may return immediately for canceled streaming requests
			next.ServeHTTP(wrapped, r)

			// Log request after handler completes (even if canceled)
			duration := time.Since(start)

			// Check if request context was canceled (e.g., client disconnected)
			canceled := ""
			if r.Context().Err() == context.Canceled {
				canceled = " [canceled]"
			}

			logger.Printf(
				"%s %s %d %s %d bytes%s",
				r.Method,
				r.URL.Path,
				wrapped.statusCode,
				duration,
				wrapped.written,
				canceled,
			)
		})
	}
}

// ----------------------------------------------------------------------------
// responseWriter wraps http.ResponseWriter to capture status code and bytes written.
// It also forwards optional interfaces like http.Flusher, http.Hijacker, etc.
// This allows the middleware to log response details without modifying handler behavior.
// ----------------------------------------------------------------------------
type responseWriter struct {
	http.ResponseWriter
	statusCode int   // Captured status code (defaults to 200)
	written    int64 // Total bytes written to response
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Default status if WriteHeader is never called
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code // Capture status code before forwarding
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n) // Track total bytes written
	return n, err
}

// Flush implements http.Flusher if the underlying ResponseWriter supports it.
// This is required for streaming responses (e.g., server-sent events).
func (rw *responseWriter) Flush() {
	if flusher, ok := rw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Hijack implements http.Hijacker if the underlying ResponseWriter supports it.
// This allows upgrading HTTP connections (e.g., to WebSocket).
func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, fmt.Errorf("http.Hijacker interface is not supported")
}

// Push implements http.Pusher if the underlying ResponseWriter supports it.
// This enables HTTP/2 server push functionality.
func (rw *responseWriter) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := rw.ResponseWriter.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return http.ErrNotSupported
}

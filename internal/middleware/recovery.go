package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

// ----------------------------------------------------------------------------
// Creates panic recovery middleware that catches panics in HTTP handlers
// and prevents them from crashing the entire server.
//
// In Go, a panic is an unrecoverable error that would normally terminate the program.
// This middleware uses Go's built-in panic recovery mechanism to catch panics,
// log them for debugging, and return a proper HTTP error response instead.
//
// Parameters:
//   - logger: A logger instance used to record panic details and stack traces
//
// Returns:
//   - A middleware function that wraps HTTP handlers with panic recovery
//
// ----------------------------------------------------------------------------
func Recovery(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// defer ensures this function runs when the surrounding function returns,
			// even if it returns due to a panic. This is the key mechanism for recovery.
			defer func() {
				// recover() catches a panic if one occurred. It returns the panic value
				// if a panic happened, or nil if execution was normal.
				// recover() only works when called directly from a deferred function.
				if err := recover(); err != nil {
					// Log the panic value and the full stack trace for debugging.
					// debug.Stack() returns the current goroutine's stack trace as a byte slice,
					// which helps identify where the panic occurred in the code.
					logger.Printf("PANIC: %v\n%s", err, debug.Stack())

					// Send a 500 Internal Server Error status to the client.
					// This provides a proper HTTP response instead of the connection
					// being abruptly closed due to the panic.
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Internal server error")
				}
			}()

			// Execute the next handler in the chain. If it panics, the deferred
			// recovery function above will catch it and handle it gracefully.
			next.ServeHTTP(w, r)
		})
	}
}

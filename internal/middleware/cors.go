package middleware

import (
	"net/http"
)

// ----------------------------------------------------------------------------
// Creates CORS middleware with the given allowed origins.
//
// CORS (Cross-Origin Resource Sharing) allows web pages from different origins
// (different protocol, domain, or port) to access this server's resources.
// ----------------------------------------------------------------------------
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	// Convert slice to map for O(1) lookup performance.
	// This is more efficient than checking a slice on every request,
	// especially when there are many allowed origins.
	originsMap := make(map[string]bool)
	for _, origin := range allowedOrigins {
		originsMap[origin] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the origin header from the incoming request.
			// The Origin header is automatically set by browsers for cross-origin requests.
			origin := r.Header.Get("Origin")

			// Check if the origin is allowed and set the appropriate CORS header.
			// If "*" is in the allowed list, all origins are permitted.
			// Otherwise, only explicitly listed origins are allowed.
			if origin != "" && (originsMap["*"] || originsMap[origin]) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			// Set CORS headers that tell the browser what is allowed:
			// - Methods: HTTP methods that can be used in cross-origin requests
			// - Headers: Request headers that can be sent (includes Connect RPC headers)
			// - Credentials: Allows cookies and authorization headers to be sent
			// - Max-Age: How long (in seconds) the browser can cache preflight responses
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Connect-Protocol-Version, Connect-Timeout-Ms, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "3600")

			// Handle preflight requests (OPTIONS method).
			// Browsers send preflight requests before certain cross-origin requests
			// to check if the actual request is allowed. We respond with 200 OK
			// and the CORS headers above, then stop processing (don't call next handler).
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// For non-preflight requests, pass control to the next handler in the chain.
			next.ServeHTTP(w, r)
		})
	}
}

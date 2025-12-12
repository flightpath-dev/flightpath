package server

import (
	"context"
	"log"
	"net/http"

	"github.com/flightpath-dev/flightpath/internal/config"
	"github.com/flightpath-dev/flightpath/internal/middleware"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// Server represents the Flightpath server. It holds the server state & provides methods to
// 1. Start and stop the server.
// 2. Register gRPC services so that incoming requests can be routed correctly.
type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
	config     *config.Config
	logger     *log.Logger
}

// NewServer creates a new Server instance
func NewServer(cfg *config.Config) *Server {
	return &Server{
		mux:    http.NewServeMux(),
		config: cfg,
		logger: log.New(log.Writer(), "[flightpath] ", log.LstdFlags|log.Lshortfile),
	}
}

// Config returns the server's configuration
func (s *Server) Config() *config.Config {
	return s.config
}

// Logger returns the server's logger
func (s *Server) Logger() *log.Logger {
	return s.logger
}

// Registers a service handler
func (s *Server) RegisterService(path string, handler http.Handler) {
	s.logger.Printf("Registering service: %s", path)
	s.mux.Handle(path, handler)
}

// Starts the HTTP server
func (s *Server) Start() error {
	addr := s.config.ServerAddr()
	handler := s.buildHandler()

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	s.logger.Printf("🚀 Flightpath server starting on %s", addr)
	s.logger.Printf("📡 Ready to accept Connect protocol requests")

	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}

// Builds the final HTTP handler with all middleware
func (s *Server) buildHandler() http.Handler {
	// Start with the mux
	handler := http.Handler(s.mux)

	// Add middleware in reverse order (last applied first)
	handler = middleware.CORS(s.config.Server.CORSOrigins)(handler)
	handler = middleware.Logging(s.logger)(handler)
	handler = middleware.Recovery(s.logger)(handler)

	// Wrap with h2c (HTTP/2 Cleartext) for Connect protocol
	// This provides HTTP/2 support for the server
	return h2c.NewHandler(handler, &http2.Server{})
}

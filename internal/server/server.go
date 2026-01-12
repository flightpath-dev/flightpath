package server

import (
	"context"
	"net/http"

	"github.com/flightpath-dev/flightpath/internal/config"
	"github.com/flightpath-dev/flightpath/internal/logger"
	"github.com/flightpath-dev/flightpath/internal/middleware"
)

// Server represents the Flightpath server. It holds the server state & provides methods to
// 1. Start and stop the server.
// 2. Register gRPC services so that incoming requests can be routed correctly.
type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
	config     *config.Config
	logger     *logger.Logger
}

// NewServer creates a new Server instance
func NewServer(cfg *config.Config) *Server {
	return &Server{
		mux:    http.NewServeMux(),
		config: cfg,
		logger: logger.New(cfg.LogLevel, cfg.LogFormat).WithPrefix("server"),
	}
}

// Config returns the server's configuration
func (s *Server) Config() *config.Config {
	return s.config
}

// Logger returns the server's logger
func (s *Server) Logger() *logger.Logger {
	return s.logger
}

// Registers a service handler
func (s *Server) RegisterService(path string, handler http.Handler) {
	s.logger.Info("Registering service", "path", path)
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

	s.logger.Info("Flightpath server starting", "address", addr)
	s.logger.Info("Ready to accept Connect protocol requests")

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

	return handler
}

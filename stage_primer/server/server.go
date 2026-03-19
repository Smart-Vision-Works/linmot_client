package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"stage_primer_config"
)

// Server represents the REST API server
type Server struct {
	router *gin.Engine
	server *http.Server
	configPath string // Allow override for testing
	// grpcHealthTarget is the address used by /healthz/grpc.
	grpcHealthTarget string
	// mockMode enables mock responses for testing without real hardware
	mockMode bool
}

// NewServer creates a new server instance
// configPath allows overriding the default config path (mainly for testing)
func NewServer(configPath ...string) (*Server, error) {
	router := gin.Default()

	srv := &Server{
		router: router,
	}

	// Set config path - default or override
	if len(configPath) > 0 {
		srv.configPath = configPath[0]
	} else {
		srv.configPath = config.DefaultConfigPath
	}

	srv.setupRoutes()

	return srv, nil
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	s.server = &http.Server{
		Addr:              addr,
		Handler:           s.router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      0,
		IdleTimeout:       120 * time.Second,
	}

	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	if s.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}

// SetGrpcHealthTarget configures the gRPC endpoint used by /healthz/grpc.
func (s *Server) SetGrpcHealthTarget(target string) {
	if s == nil {
		return
	}
	s.grpcHealthTarget = target
}

// SetMockMode enables or disables mock mode for testing.
// When enabled, USB device discovery returns fake devices.
func (s *Server) SetMockMode(enabled bool) {
	if s == nil {
		return
	}
	s.mockMode = enabled
}

// Handler returns the HTTP handler for testing purposes
func (s *Server) Handler() http.Handler {
	return s.router
}

package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// setupRoutes configures the API routes
func (s *Server) setupRoutes() {
	// Keep parity with other SVW services that expose project metadata used by
	// setup verification (/data/project.txt, /data/sha.txt, etc).
	s.router.Static("/data", "/opt/public/data")

	s.router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	s.router.GET("/healthz/grpc", func(c *gin.Context) {
		target := s.grpcHealthTarget
		if target == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unavailable",
				"error":  "grpc target not configured",
			})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
		defer cancel()

		conn, err := grpc.DialContext(ctx, target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unavailable",
				"error":  err.Error(),
			})
			return
		}
		defer conn.Close()

		client := grpc_health_v1.NewHealthClient(conn)
		resp, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
		if err != nil || resp.GetStatus() != grpc_health_v1.HealthCheckResponse_SERVING {
			if err == nil {
				err = errors.New("grpc not serving")
			}
			status := "unavailable"
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": status,
				"error":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

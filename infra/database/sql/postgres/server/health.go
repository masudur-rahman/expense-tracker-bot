package server

import (
	"context"

	"github.com/masudur-rahman/expense-tracker-bot/infra/logr"

	health "google.golang.org/grpc/health/grpc_health_v1"
)

type HealthChecker struct {
	isDatabaseReady bool
}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{}
}

func (s *HealthChecker) setDatabaseReady() {
	s.isDatabaseReady = true
}

func (s *HealthChecker) Check(ctx context.Context, req *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	logr.DefaultLogger.Infow("Serving the Check request for health check")
	status := health.HealthCheckResponse_NOT_SERVING
	if s.isDatabaseReady {
		status = health.HealthCheckResponse_SERVING
	}

	return &health.HealthCheckResponse{
		Status: status,
	}, nil
}

func (s *HealthChecker) Watch(req *health.HealthCheckRequest, server health.Health_WatchServer) error {
	logr.DefaultLogger.Infow("Serving the Watch request for health check")
	status := health.HealthCheckResponse_NOT_SERVING
	if s.isDatabaseReady {
		status = health.HealthCheckResponse_SERVING
	}

	return server.Send(&health.HealthCheckResponse{
		Status: status,
	})
}

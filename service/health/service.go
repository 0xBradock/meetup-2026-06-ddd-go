package health

import (
	"context"
	"log/slog"
)

// Service is the application service for the health domain. It orchestrates
// domain logic and delegates persistence to the injected Repository.
type Service struct {
	repository Repository
	logger     *slog.Logger
}

// NewService creates a Service backed by the given Repository.
// Pass a logger enriched with domain-specific attributes (e.g. logger.With("domain", "health")).
// Pass slog.New(slog.NewTextHandler(io.Discard, nil)) in tests to suppress output.
func NewService(repository Repository, logger *slog.Logger) *Service {
	return &Service{repository: repository, logger: logger}
}

// GetHealthOutput is the result of a GetHealth call.
type GetHealthOutput struct {
	Status        string `json:"status"`
	CheckedAtUnix int64  `json:"checkedAtUnix"`
}

// GetHealth returns the current health status of the service.
func (s *Service) GetHealth(ctx context.Context) (GetHealthOutput, error) {
	s.logger.DebugContext(ctx, "getting health status")

	health, err := s.repository.GetHealth(ctx)
	if err != nil {
		return GetHealthOutput{}, err
	}

	return GetHealthOutput(health), nil
}

// PingInput carries the input for a Ping call.
type PingInput struct {
	Message string
}

// PingOutput is the result of a Ping call.
type PingOutput struct {
	Message        string `json:"message"`
	ReceivedAtUnix int64  `json:"receivedAtUnix"`
}

// Ping stores a ping message and returns it with a received timestamp.
func (s *Service) Ping(ctx context.Context, input PingInput) (PingOutput, error) {
	ping, err := s.repository.Ping(ctx, input.Message)
	if err != nil {
		return PingOutput{}, err
	}

	s.logger.DebugContext(ctx, "ping stored")

	return PingOutput(ping), nil
}

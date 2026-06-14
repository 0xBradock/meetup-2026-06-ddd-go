package health

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestServicePing(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{}
	svc := NewService(repo, discardLogger())

	_, err := svc.Ping(context.Background(), PingInput{Message: "  "})
	if !errors.Is(err, ErrEmptyPingMessage) {
		t.Fatalf("expected ErrEmptyPingMessage, got %v", err)
	}
}

func TestServiceGetHealth(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{
		health: Health{Status: "ok", CheckedAtUnix: 42},
	}
	svc := NewService(repo, discardLogger())

	out, err := svc.GetHealth(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.Status != "ok" || out.CheckedAtUnix != 42 {
		t.Fatalf("unexpected output: %+v", out)
	}
}

type fakeRepository struct {
	health Health
	ping   Ping
	err    error
}

func (f *fakeRepository) GetHealth(context.Context) (Health, error) {
	if f.err != nil {
		return Health{}, f.err
	}
	return f.health, nil
}

func (f *fakeRepository) Ping(_ context.Context, message string) (Ping, error) {
	if f.err != nil {
		return Ping{}, f.err
	}
	ping, err := NewPing(message)
	if err != nil {
		return Ping{}, err
	}
	f.ping = ping
	return ping, nil
}

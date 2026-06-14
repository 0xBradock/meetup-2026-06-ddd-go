package inmem

import (
	"context"
	"testing"
)

func TestRepositoryPing(t *testing.T) {
	t.Parallel()

	repo := NewRepository()

	ping, err := repo.Ping(context.Background(), "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ping.Message != "hello" {
		t.Fatalf("unexpected ping message: %q", ping.Message)
	}
}

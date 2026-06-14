package db

import (
	"context"
	"testing"
)

func TestNewPG(t *testing.T) {
	t.Run("non-postgres driver returns in-memory repositories", func(t *testing.T) {
		t.Parallel()

		repos, cleanup, err := NewPG(context.Background(), "inmemory", func(string) string { return "" })
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer cleanup()

		h, err := repos.Health.GetHealth(context.Background())
		if err != nil {
			t.Fatalf("unexpected health error: %v", err)
		}
		if h.Status == "" {
			t.Fatal("expected non-empty health status")
		}
	})

	t.Run("postgres driver requires DATABASE_URL", func(t *testing.T) {
		t.Parallel()

		_, _, err := NewPG(context.Background(), "postgres", func(string) string { return "" })
		if err == nil {
			t.Fatal("expected error when DATABASE_URL is missing")
		}
	})
}

func TestNewMSSQL(t *testing.T) {
	t.Run("non-mssql driver returns in-memory repositories", func(t *testing.T) {
		t.Parallel()

		repos, cleanup, err := NewMSSQL(context.Background(), "inmemory", func(string) string { return "" })
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer cleanup()

		h, err := repos.Health.GetHealth(context.Background())
		if err != nil {
			t.Fatalf("unexpected health error: %v", err)
		}
		if h.Status == "" {
			t.Fatal("expected non-empty health status")
		}
	})

	t.Run("mssql driver requires MSSQL_DATABASE_URL", func(t *testing.T) {
		t.Parallel()

		_, _, err := NewMSSQL(context.Background(), "mssql", func(string) string { return "" })
		if err == nil {
			t.Fatal("expected error when MSSQL_DATABASE_URL is missing")
		}
	})
}

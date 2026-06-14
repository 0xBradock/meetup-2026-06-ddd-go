package inmem_test

import (
	"context"
	"testing"
	"time"

	"go-svr/db/inmem"
)

func TestCache_GetMissReturnsEmpty(t *testing.T) {
	t.Parallel()
	c := inmem.NewCache()
	val, err := c.Get(context.Background(), "missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "" {
		t.Fatalf("expected empty string on miss, got %q", val)
	}
}

func TestCache_SetGet(t *testing.T) {
	t.Parallel()
	c := inmem.NewCache()
	ctx := context.Background()

	if err := c.Set(ctx, "k", "v", time.Minute); err != nil {
		t.Fatalf("Set: %v", err)
	}
	val, err := c.Get(ctx, "k")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if val != "v" {
		t.Fatalf("expected %q, got %q", "v", val)
	}
}

func TestCache_ExpiredEntryTreatedAsMiss(t *testing.T) {
	t.Parallel()
	c := inmem.NewCache()
	ctx := context.Background()

	if err := c.Set(ctx, "k", "v", -time.Second); err != nil {
		t.Fatalf("Set: %v", err)
	}
	val, err := c.Get(ctx, "k")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if val != "" {
		t.Fatalf("expected empty string for expired entry, got %q", val)
	}
}

func TestCache_Del(t *testing.T) {
	t.Parallel()
	c := inmem.NewCache()
	ctx := context.Background()

	_ = c.Set(ctx, "a", "1", time.Minute)
	_ = c.Set(ctx, "b", "2", time.Minute)

	if err := c.Del(ctx, "a", "b"); err != nil {
		t.Fatalf("Del: %v", err)
	}
	for _, k := range []string{"a", "b"} {
		val, _ := c.Get(ctx, k)
		if val != "" {
			t.Fatalf("expected key %q to be deleted", k)
		}
	}
}

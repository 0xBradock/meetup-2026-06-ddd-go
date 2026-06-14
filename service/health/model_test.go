package health

import "testing"

func TestNewPing(t *testing.T) {
	t.Parallel()

	t.Run("rejects empty message", func(t *testing.T) {
		t.Parallel()

		_, err := NewPing("   ")
		if err == nil {
			t.Fatal("expected error for empty ping message")
		}
	})

	t.Run("creates ping with message", func(t *testing.T) {
		t.Parallel()

		ping, err := NewPing("hello")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if ping.Message != "hello" {
			t.Fatalf("unexpected message: %q", ping.Message)
		}

		if ping.ReceivedAtUnix == 0 {
			t.Fatal("expected timestamp to be set")
		}
	})
}

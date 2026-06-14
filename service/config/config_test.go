package config

import (
	"os"
	"testing"
	"time"
)

func setBaseServerEnv(t *testing.T) {
	t.Helper()
	t.Setenv("PG_DRIVER", "inmemory")
	t.Setenv("MSSQL_DRIVER", "inmemory")
	t.Setenv("SERVICE_NAME", "")
	t.Setenv("DATABASE_URL", "")
	t.Setenv("ENVIRONMENT", "")
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("HTTP_READ_TIMEOUT", "")
	t.Setenv("HTTP_WRITE_TIMEOUT", "")
	t.Setenv("HTTP_IDLE_TIMEOUT", "")
	t.Setenv("SERVER_SHUTDOWN_TIME", "")
}

func TestLoadServerConfig_ServiceName(t *testing.T) {
	setBaseServerEnv(t)

	cfg, err := LoadServerConfig(os.Getenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.ServiceName != SERVICE_NAME {
		t.Fatalf("expected default service name %q, got %q", SERVICE_NAME, cfg.ServiceName)
	}

	t.Setenv("SERVICE_NAME", "something-else")

	cfg, err = LoadServerConfig(os.Getenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.ServiceName != "something-else" {
		t.Fatalf("expected overridden service name %q, got %q", "something-else", cfg.ServiceName)
	}
}

func TestLoadServerConfig_Environment(t *testing.T) {
	t.Run("uses default when empty", func(t *testing.T) {
		setBaseServerEnv(t)

		cfg, err := LoadServerConfig(os.Getenv)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Environment != "recette" {
			t.Fatalf("expected default environment recette, got %q", cfg.Environment)
		}
	})

	t.Run("accepts allowed environments", func(t *testing.T) {
		cases := []string{"prod", "preprod", "recette"}

		for _, env := range cases {
			env := env
			t.Run(env, func(t *testing.T) {
				setBaseServerEnv(t)
				t.Setenv("ENVIRONMENT", env)

				cfg, err := LoadServerConfig(os.Getenv)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if cfg.Environment != env {
					t.Fatalf("expected environment %q, got %q", env, cfg.Environment)
				}
			})
		}
	})

	t.Run("rejects unsupported environment", func(t *testing.T) {
		setBaseServerEnv(t)
		t.Setenv("ENVIRONMENT", "dev")

		_, err := LoadServerConfig(os.Getenv)
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestLoadServerConfig_DefaultTimeouts(t *testing.T) {
	setBaseServerEnv(t)

	cfg, err := LoadServerConfig(os.Getenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.HTTPReadTimeout != 5*time.Second {
		t.Fatalf("expected HTTP read timeout 5s, got %s", cfg.HTTPReadTimeout)
	}
	if cfg.HTTPWriteTimeout != 5*time.Second {
		t.Fatalf("expected HTTP write timeout 5s, got %s", cfg.HTTPWriteTimeout)
	}
	if cfg.HTTPIdleTimeout != 5*time.Second {
		t.Fatalf("expected HTTP idle timeout 5s, got %s", cfg.HTTPIdleTimeout)
	}
}

func TestLoadServerConfig_InvalidTimeout(t *testing.T) {
	setBaseServerEnv(t)
	t.Setenv("HTTP_READ_TIMEOUT", "nope")

	_, err := LoadServerConfig(os.Getenv)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadServerConfig_ServerShutdownTime(t *testing.T) {
	t.Run("uses default when empty", func(t *testing.T) {
		setBaseServerEnv(t)

		cfg, err := LoadServerConfig(os.Getenv)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.ServerShutdownTime != 5 {
			t.Fatalf("expected default shutdown time 5, got %d", cfg.ServerShutdownTime)
		}
	})

	t.Run("parses value", func(t *testing.T) {
		setBaseServerEnv(t)
		t.Setenv("SERVER_SHUTDOWN_TIME", "12")

		cfg, err := LoadServerConfig(os.Getenv)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.ServerShutdownTime != 12 {
			t.Fatalf("expected shutdown time 12, got %d", cfg.ServerShutdownTime)
		}
	})

	t.Run("rejects invalid value", func(t *testing.T) {
		setBaseServerEnv(t)
		t.Setenv("SERVER_SHUTDOWN_TIME", "abc")

		_, err := LoadServerConfig(os.Getenv)
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

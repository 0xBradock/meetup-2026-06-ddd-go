package main

import (
	"bytes"
	"context"
	"errors"
	"net"
	"os"
	"strings"
	"testing"

	"go-svr/config"
)

func setBaseRunEnv(t *testing.T) {
	t.Helper()
	t.Setenv("PG_DRIVER", "inmemory")
	t.Setenv("MSSQL_DRIVER", "inmemory")
	t.Setenv("SERVICE_NAME", "")
	t.Setenv("DATABASE_URL", "")
	t.Setenv("ENVIRONMENT", "")
	t.Setenv("HTTP_ADDR", ":0")
	t.Setenv("HTTP_READ_TIMEOUT", "")
	t.Setenv("HTTP_WRITE_TIMEOUT", "")
	t.Setenv("HTTP_IDLE_TIMEOUT", "")
	t.Setenv("SERVER_SHUTDOWN_TIME", "")
	t.Setenv("CACHE_URL", "")
	t.Setenv("CACHE_TLS_CA_CERT", "")
	t.Setenv("CACHE_PREFIX", "")
	t.Setenv("CACHE_TTL", "")
	t.Setenv("CACHE_HEALTH_TTL", "")
}

func TestRun_LogsIncludeServiceAndEnvironment(t *testing.T) {
	setBaseRunEnv(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	stdout := &bytes.Buffer{}
	err := Run(ctx, nil, os.Getenv, stdout, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	logs := stdout.String()
	if !strings.Contains(logs, `"service":"`+config.SERVICE_NAME+`"`) {
		t.Fatalf("expected logs to include service field, got %s", logs)
	}
	if !strings.Contains(logs, `"environment":"recette"`) {
		t.Fatalf("expected logs to include environment field, got %s", logs)
	}
}

func TestRun_LogsIncludeOverriddenServiceName(t *testing.T) {
	setBaseRunEnv(t)
	t.Setenv("SERVICE_NAME", "orders-api")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	stdout := &bytes.Buffer{}
	err := Run(ctx, nil, os.Getenv, stdout, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	logs := stdout.String()
	if !strings.Contains(logs, `"service":"orders-api"`) {
		t.Fatalf("expected logs to include overridden service field, got %s", logs)
	}
}

func TestRun_InvalidDriver(t *testing.T) {
	setBaseRunEnv(t)
	t.Setenv("PG_DRIVER", "invalid")

	err := Run(context.Background(), nil, os.Getenv, &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_ListenerFailure(t *testing.T) {
	setBaseRunEnv(t)

	originalListen := listenFn
	listenFn = func(string, string) (net.Listener, error) {
		return nil, errors.New("listen failed")
	}
	t.Cleanup(func() { listenFn = originalListen })

	err := Run(context.Background(), nil, os.Getenv, &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_ContextCancelled(t *testing.T) {
	setBaseRunEnv(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Run(ctx, nil, os.Getenv, &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

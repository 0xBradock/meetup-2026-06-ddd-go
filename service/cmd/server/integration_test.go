//go:build integration

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

type scenarioCtx struct {
	baseURL  string
	client   *http.Client
	response *http.Response
	body     []byte
}

func (s *scenarioCtx) sendGET(path string) error {
	resp, err := s.client.Get(s.baseURL + path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	s.response = resp
	s.body, err = io.ReadAll(resp.Body)
	return err
}

func (s *scenarioCtx) sendPOST(path string, body *godog.DocString) error {
	resp, err := s.client.Post(s.baseURL+path, "application/json", strings.NewReader(body.Content))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	s.response = resp
	s.body, err = io.ReadAll(resp.Body)
	return err
}

func (s *scenarioCtx) statusShouldBe(expected int) error {
	if s.response.StatusCode != expected {
		return fmt.Errorf("expected status %d, got %d (body: %s)", expected, s.response.StatusCode, s.body)
	}
	return nil
}

func (s *scenarioCtx) bodyHasField(field string) error {
	var m map[string]any
	if err := json.Unmarshal(s.body, &m); err != nil {
		return fmt.Errorf("failed to parse body as JSON: %w", err)
	}
	if _, ok := m[field]; !ok {
		return fmt.Errorf("field %q not found in body: %s", field, s.body)
	}
	return nil
}

func (s *scenarioCtx) bodyHasMessage(msg string) error {
	var m map[string]any
	if err := json.Unmarshal(s.body, &m); err != nil {
		return fmt.Errorf("failed to parse body as JSON: %w", err)
	}
	if m["message"] != msg {
		return fmt.Errorf("expected message %q, got %v (body: %s)", msg, m["message"], s.body)
	}
	return nil
}

func TestIntegration(t *testing.T) {
	baseURL, shutdown := startTestServer(t)
	t.Cleanup(shutdown)

	suite := godog.TestSuite{
		ScenarioInitializer: func(sc *godog.ScenarioContext) {
			ctx := &scenarioCtx{
				baseURL: baseURL,
				client:  &http.Client{Timeout: 5 * time.Second},
			}
			sc.Step(`^I send a GET request to "([^"]*)"$`, ctx.sendGET)
			sc.Step(`^I send a POST request to "([^"]*)" with body:$`, ctx.sendPOST)
			sc.Step(`^the response status should be (\d+)$`, ctx.statusShouldBe)
			sc.Step(`^the response body should contain a "([^"]*)" field$`, ctx.bodyHasField)
			sc.Step(`^the response body should contain message "([^"]*)"$`, ctx.bodyHasMessage)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../tests/features"},
			TestingT: t,
			Strict:   true, // undefined or pending steps fail the suite
		},
	}

	if suite.Run() != 0 {
		t.Fatal("integration test suite failed")
	}
}

func startTestServer(t *testing.T) (baseURL string, shutdown func()) {
	t.Helper()

	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("create listener: %v", err)
	}
	addr := lis.Addr().String()

	orig := listenFn
	listenFn = func(_, _ string) (net.Listener, error) { return lis, nil }
	t.Cleanup(func() { listenFn = orig })

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)

	go func() {
		done <- Run(ctx, nil, func(k string) string {
			switch k {
			case "PG_DRIVER", "MSSQL_DRIVER":
				return "inmemory"
			case "HTTP_ADDR":
				return ":0"
			}
			return ""
		}, io.Discard, io.Discard)
	}()

	base := "http://" + addr
	if err := pollReady(base+"/health", 3*time.Second); err != nil {
		cancel()
		t.Fatalf("server not ready: %v", err)
	}

	return base, func() {
		cancel()
		<-done
	}
}

func pollReady(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url) //nolint:noctx
		if err == nil {
			resp.Body.Close()
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	return fmt.Errorf("timeout after %s waiting for %s", timeout, url)
}

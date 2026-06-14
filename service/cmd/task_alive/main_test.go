package main

import (
	"context"
	"io"
	"testing"
)

func TestRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		env     map[string]string
		wantErr bool
	}{
		{
			name:    "missing TASK_ID returns error",
			env:     map[string]string{"PG_DRIVER": "inmemory"},
			wantErr: true,
		},
		{
			name: "valid TASK_ID with inmemory DB succeeds",
			env: map[string]string{
				"PG_DRIVER": "inmemory",
				"TASK_ID":   "test-task-id-123",
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			getenv := func(key string) string { return tc.env[key] }
			err := run(context.Background(), nil, getenv, io.Discard, io.Discard)
			if (err != nil) != tc.wantErr {
				t.Fatalf("run() error = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

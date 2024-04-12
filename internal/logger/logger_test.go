package logger

import (
	"context"
	"testing"
)

func TestLoggerFromContext(t *testing.T) {
	log := FromContext(context.Background())
	if log == nil {
		t.Fatalf("want slog.Default() got nil")
	}
}

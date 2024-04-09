package testutils

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/holos-run/holos/internal/server/app"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	"github.com/lmittmann/tint"
)

type testWriter struct {
	t testing.TB
}

func (w *testWriter) Write(b []byte) (int, error) {
	w.t.Logf("%s", b)
	return len(b), nil
}

// TestLogger is an adapter that sends output to t.Log so that log messages are
// associated with failing tests.
func TestLogger(t testing.TB) *slog.Logger {
	t.Helper()
	testHandler := tint.NewHandler(&testWriter{t}, &tint.Options{
		Level:     slog.LevelDebug,
		AddSource: true,
		NoColor:   true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if len(groups) == 0 {
				switch key := a.Key; key {
				case slog.TimeKey:
					return slog.Attr{} // Remove the timestamp
				case slog.SourceKey:
					if src, ok := a.Value.Any().(*slog.Source); ok {
						name := fmt.Sprintf("%s:%d:", filepath.Base(src.File), src.Line)
						return slog.String("source", name)
					}
				}
			}
			return a
		},
	})
	return slog.New(testHandler)
}

// LogCtx returns a new background context.Context carrying a *slog.Logger wired
// to t.Log(). Useful to associate all logs with the test case that caused them.
func LogCtx(t testing.TB) context.Context {
	return logger.NewContext(context.Background(), TestLogger(t))
}

// NewAppContext returns a new app.App wired to t.Log().
func NewAppContext(t testing.TB) app.App {
	log := TestLogger(t)
	return app.App{
		Context: logger.NewContext(context.Background(), log),
		Logger:  log,
	}
}

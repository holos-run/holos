// Package logger logs http responses
// See: https://github.com/elithrar/admission-control/blob/v0.6.7/request_logger.go#L40
package logger

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/holos-run/holos/pkg/logger"
)

func NewContext(ctx context.Context, log *slog.Logger) context.Context {
	return logger.NewContext(ctx, log)
}

// FromContext returns the *slog.Logger previously stored in ctx by NewContext.
// slog.Default() is returned otherwise.
func FromContext(ctx context.Context) *slog.Logger {
	return logger.FromContext(ctx)
}

// LoggingMiddleware returns a handler that adds a *slog.Logger to the request
// context.Context retrievable by FromContext. The returned handler is useful as
// the outer client facing edge of a middleware chain and includes attributes on
// the log messages.
func LoggingMiddleware(log *slog.Logger) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			// Test cases inject a logger wired to t.Log(), use it if present.
			if logContext := logger.FromContextMaybe(r.Context()); logContext != nil {
				log = logContext
			}

			log = log.With(
				"proto", r.Proto,
				"uri", r.URL.RequestURI(),
				"method", r.Method,
				"remote", GetClientIP(r),
				"user-agent", r.UserAgent(),
			)

			ctx := NewContext(r.Context(), log)
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r.WithContext(ctx))
			log.DebugContext(ctx, "response", "code", wrapped.code(), "duration", time.Since(start))
		}
		return http.HandlerFunc(fn)
	}
}

// GetClientIP returns the client address from the x-forwarded-for header or the
// http.Request RemoteAddr field.
func GetClientIP(r *http.Request) string {
	// Try to get the client IP from the X-Forwarded-For header.
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// The X-Forwarded-For header can be a comma-separated list of IPs.
		// The client's IP is typically the first one.
		client, _, _ := strings.Cut(xff, ",")
		return strings.TrimSpace(client)
	}

	return r.RemoteAddr
}

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) code() int {
	if rw.status == 0 {
		return http.StatusOK
	}
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

// Package logger logs http responses
// See: https://github.com/elithrar/admission-control/blob/v0.6.7/request_logger.go#L40
package logger

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// key is an unexported type for keys defined in this package to prevent
// collisions with keys defined in other packages.
type key int

// https://cs.opensource.google/go/go/+/refs/tags/go1.21.1:src/context/context.go;l=140-158
// loggerKey is the key for *slog.Logger values in Contexts. It us unexported;
// clients use NewContext and FromContext instead of this key directly.
var loggerKey key

// NewContext returns a new Context that carries value logger. Use FromContext
// to retrieve the value.
func NewContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext returns the *slog.Logger previously stored in ctx by NewContext.
// slog.Default() is returned otherwise.
func FromContext(ctx context.Context) *slog.Logger {
	// https://go.dev/ref/spec#Type_assertions
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// LoggingMiddleware returns a handler that adds a *slog.Logger to the request
// context.Context retrievable by FromContext. The returned handler is useful as
// the outer client facing edge of a middleware chain and includes attributes on
// the log messages.
func LoggingMiddleware(logger *slog.Logger) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := logger
			start := time.Now()
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			// Test cases inject a logger wired to t.Log(), use it if present.
			if loggerFromContext, ok := r.Context().Value(loggerKey).(*slog.Logger); ok {
				log = loggerFromContext
			}

			log = log.With(
				"proto", r.Proto,
				"uri", r.URL.RequestURI(),
				"method", r.Method,
				"remote", GetClientIP(r),
				"user-agent", r.UserAgent(),
			)

			ctx := NewContext(r.Context(), logger)
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

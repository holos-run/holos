package signals

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

type Shutdown struct {
	logger                *slog.Logger
	serverShutdownTimeout time.Duration
}

func NewShutdown(serverShutdownTimeout time.Duration, logger *slog.Logger) *Shutdown {
	srv := &Shutdown{
		logger:                logger,
		serverShutdownTimeout: serverShutdownTimeout,
	}

	return srv
}

func (s *Shutdown) Graceful(stopCh <-chan struct{}, httpServer *http.Server, healthy *int32, ready *int32) {
	ctx := context.Background()

	// wait for SIGTERM or SIGINT
	<-stopCh
	ctx, cancel := context.WithTimeout(ctx, s.serverShutdownTimeout)
	defer cancel()

	// all calls to /healthz and /readyz will fail from now on
	atomic.StoreInt32(healthy, 0)
	atomic.StoreInt32(ready, 0)

	delay := 9 * time.Second
	delayOverride := os.Getenv("SHUTDOWN_DELAY")
	if delayOverride != "" {
		if val, err := strconv.Atoi(delayOverride); err != nil {
			s.logger.ErrorContext(ctx, "could not override delay, SHUTDOWN_DELAY env val is not an int", "delay", delay, "err", err)
		} else {
			delay = time.Duration(val) * time.Second
			s.logger.DebugContext(ctx, "SHUTDOWN_DELAY env override in effect", "delay", delay)
		}
	}
	s.logger.DebugContext(ctx, "shutting down http/https server", "delay", delay, "timeout", s.serverShutdownTimeout)
	// wait for Kubernetes readiness probe to remove this instance from the load balancer
	// the readiness check interval must be lower than the timeout
	time.Sleep(delay)

	// determine if the http server was started
	if httpServer != nil {
		if err := httpServer.Shutdown(ctx); err != nil {
			s.logger.ErrorContext(ctx, "could not shutdown http server gracefully", "err", err)
		}
	}
}

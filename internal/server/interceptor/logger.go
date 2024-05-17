package interceptor

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/holos-run/holos/internal/server/middleware/logger"
)

func NewLogger() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()
			rpcLogger := logger.FromContext(ctx).With(
				"procedure", req.Spec().Procedure,
				"request_id", req.Header().Get("x-request-id"),
				"traceid", req.Header().Get("x-b3-traceid"),
				"spanid", req.Header().Get("x-b3-spanid"),
				"parentspanid", req.Header().Get("x-b3-parentspanid"),
				"sampled", req.Header().Get("x-b3-sampled"),
				"host", req.Header().Get("host"),
			)
			ctx = logger.NewContext(ctx, rpcLogger)
			resp, err := next(ctx, req)
			go emitLog(ctx, start, err)
			return resp, err
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func emitLog(ctx context.Context, start time.Time, err error) {
	log := logger.FromContext(ctx)
	if err == nil {
		log = log.With("ok", true)
	} else {
		log = log.With("ok", false, "code", connect.CodeOf(err), "err", err)
	}
	log.InfoContext(ctx, "response", "duration", time.Since(start))
}

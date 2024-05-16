package cli

import (
	"context"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	cue "cuelang.org/go/cue/errors"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

// MakeMain makes a main function for the cli or tests.
func MakeMain(options ...holos.Option) func() int {
	return func() (exitCode int) {
		cfg := holos.New(options...)
		slog.SetDefault(cfg.Logger())
		ctx := context.Background()
		if err := New(cfg).ExecuteContext(ctx); err != nil {
			return HandleError(ctx, err, cfg)
		}
		return 0
	}
}

// HandleError is the top level error handler that unwraps and logs errors.
func HandleError(ctx context.Context, err error, hc *holos.Config) (exitCode int) {
	// Connect errors have codes, log them.
	log := hc.NewTopLevelLogger().With("code", connect.CodeOf(err))
	var cueErr cue.Error
	var errAt *errors.ErrorAt
	const msg = "could not execute"
	if errors.As(err, &errAt) {
		log.ErrorContext(ctx, msg, "err", errAt.Unwrap(), "loc", errAt.Source.Loc())
	} else {
		log.ErrorContext(ctx, msg, "err", err)
	}
	// cue errors are bundled up as a list and refer to multiple files / lines.
	if errors.As(err, &cueErr) {
		msg := cue.Details(cueErr, nil)
		_, _ = fmt.Fprint(hc.Stderr(), msg)
	}
	// connect errors have details and codes.
	// Refer to https://connectrpc.com/docs/go/errors
	if connectErr := new(connect.Error); errors.As(err, &connectErr) {
		for _, detail := range connectErr.Details() {
			msg, valueErr := detail.Value()
			if valueErr != nil {
				log.WarnContext(ctx, "could not decode error detail", "err", err, "type", detail.Type(), "note", "this usually means we don't have the schema for the protobuf message type")
				continue
			}
			if info, ok := msg.(*errdetails.ErrorInfo); ok {
				logDetail := log.With("reason", info.GetReason(), "domain", info.GetDomain())
				for k, v := range info.GetMetadata() {
					logDetail = logDetail.With(k, v)
				}
				logDetail.ErrorContext(ctx, info.String())
			}
		}
	}

	return 1
}

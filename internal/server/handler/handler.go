// Package handler implements service handlers. Handlers are not bound to a
// specific server transport like http or nats, instead they take a context and
// a request proto buf as inputs. http and other transports.
//
// New handler constructors are expected to take in dependencies provided by
// other packages like connected database clients, loggers, tracers, etc...
package handler

import (
	"context"
	"fmt"

	"github.com/holos-run/holos/internal/ent"
	"github.com/holos-run/holos/internal/server/middleware/logger"
)

// WithTx runs callbacks in a transaction as described in https://entgo.io/docs/transactions/#best-practices
func WithTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
	log := logger.FromContext(ctx)
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			log.ErrorContext(ctx, "panic", "v", v)
			_ = tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			log.ErrorContext(ctx, "could not roll back tx", "err", rerr)
			err = fmt.Errorf("coult not roll back tx: %w: %w", rerr, err)
		} else {
			log.WarnContext(ctx, "rolled back failed tx", "err", err)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		log.ErrorContext(ctx, "could not commit transaction", "err", err)
		return fmt.Errorf("could not commit: %w", err)
	}
	return nil
}

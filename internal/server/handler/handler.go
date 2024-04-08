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
	"log/slog"

	"github.com/holos-run/holos/internal/server/ent"
)

// WithTx runs callbacks in a transaction as described in https://entgo.io/docs/transactions/#best-practices
func WithTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			slog.ErrorContext(ctx, "panic", "v", v)
			_ = tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

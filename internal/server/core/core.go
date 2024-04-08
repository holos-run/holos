// Package core contains essential structs and runtime environment structures for the application.

package core

import (
	"context"
	"log/slog"
)

// AppContext carries core application level elements commonly passed around the
// entire system. The primary use case is to inject a context.Context and a
// slog.Logger.
type AppContext struct {
	Context context.Context
	Logger  *slog.Logger
	Config  *Config
}

// ContextLogger is a convenience method to easily unpack the passed AppContext
// into local variables.
//
//	ctx, log := app.ContextLogger()
func (a AppContext) ContextLogger() (ctx context.Context, log *slog.Logger) {
	ctx = a.Context
	if ctx == nil {
		ctx = context.Background()
	}
	log = a.Logger
	if log == nil {
		log = slog.Default()
		log.WarnContext(ctx, "programming error: ensure AppContext has a logger")
	}
	return
}

func (a AppContext) WithContext(ctx context.Context) AppContext {
	a.Context = ctx
	return a
}

func (a AppContext) WithLogger(logger *slog.Logger) AppContext {
	a.Logger = logger
	return a
}

// NewAppContext returns a new AppContext initialized with default values.
func NewAppContext() AppContext {
	return AppContext{
		Context: context.Background(),
		Logger:  slog.Default(),
		Config:  NewConfig(),
	}
}

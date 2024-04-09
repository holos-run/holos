package app

import (
	"context"
	"log/slog"
)

// App carries core application level elements commonly passed around the
// entire system. The primary use case is to inject a context.Context and a
// slog.Logger.
type App struct {
	Context context.Context
	Logger  *slog.Logger
	Config  *Config
}

// ContextLogger is a convenience method to easily unpack the passed AppContext
// into local variables.
//
//	ctx, log := app.ContextLogger()
func (a App) ContextLogger() (ctx context.Context, log *slog.Logger) {
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

func (a App) WithContext(ctx context.Context) App {
	a.Context = ctx
	return a
}

func (a App) WithLogger(logger *slog.Logger) App {
	a.Logger = logger
	return a
}

// New returns a new App initialized with default values.
func New() App {
	return App{
		Context: context.Background(),
		Logger:  slog.Default(),
		Config:  NewConfig(),
	}
}

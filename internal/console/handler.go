package console

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
)

func NewHandler(out io.Writer, opts *Options) slog.Handler {
	h := &ConsoleHandler{out: out, mu: &sync.Mutex{}}
	if opts != nil {
		h.opts = *opts
	}
	if h.opts.Level == nil {
		h.opts.Level = slog.LevelInfo
	}
	return h
}

type Options struct {
	Level slog.Leveler
}

type ConsoleHandler struct {
	opts Options
	mu   *sync.Mutex
	out  io.Writer
}

func (h *ConsoleHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *ConsoleHandler) Handle(ctx context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := fmt.Fprintln(h.out, r.Message)
	return err
}

func (h *ConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *ConsoleHandler) WithGroup(name string) slog.Handler {
	return h
}

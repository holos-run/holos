package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
)

// Source represents the Source file and line where an error was first
// encountered by holos code.
type Source struct {
	File string `json:"file"`
	Line int    `json:"line"`
}

func (s *Source) Loc() string {
	return fmt.Sprintf("%s:%d", filepath.Base(s.File), s.Line)
}

// HolosError wraps an error with the Source location the error was encountered
// within Holos for tracing from the top level error handler.
type HolosError struct {
	Err    error
	Source Source
}

// Unwrap implements error wrapping.
func (e *HolosError) Unwrap() error {
	return e.Err
}

// Error returns the error string with Source location.
func (e *HolosError) Error() string {
	return e.Source.Loc() + ": " + e.Err.Error()
}

// WrapError wraps err in a HolosError or returns err if err is nil, already a
// HolosError, or caller info is not available.
func WrapError(err error) error {
	// Nothing to do
	if err == nil {
		return nil
	}

	// Already a holos error no need to do anything.
	var herr *HolosError
	if errors.As(err, &herr) {
		return err
	}

	// Try to wrap err with caller info
	if _, file, line, ok := runtime.Caller(1); ok {
		return &HolosError{
			Err: err,
			Source: Source{
				File: file,
				Line: line,
			},
		}
	}

	return err
}

// LogError logs Err with Source location if Err is a HolosError
func LogError(log *slog.Logger, ctx context.Context, level slog.Level, err error, msg string, args ...any) {
	var herr *HolosError
	if ok := errors.As(err, &herr); ok {
		args = append(args,
			slog.String("err", herr.Unwrap().Error()),
			slog.String("loc", herr.Source.Loc()),
		)
	} else {
		if _, file, line, ok := runtime.Caller(1); ok {
			source := Source{file, line}
			args = append(args, slog.String("loc", source.Loc()))
		}
		args = append(args, slog.String("err", err.Error()))
	}
	log.Log(ctx, level, msg, args...)
}

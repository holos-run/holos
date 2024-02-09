// Package wrapper provides error wrapping with location information
package wrapper

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
)

// Source represents the Source file and line where an error was encountered
type Source struct {
	File string `json:"file"`
	Line int    `json:"line"`
}

func (s *Source) Loc() string {
	return fmt.Sprintf("%s:%d", filepath.Base(s.File), s.Line)
}

// ErrorAt wraps an error with the Source location the error was encountered at
// for tracing from a top level error handler.
type ErrorAt struct {
	Err    error
	Source Source
}

// Unwrap implements error wrapping.
func (e *ErrorAt) Unwrap() error {
	return e.Err
}

// Error returns the error string with Source location.
func (e *ErrorAt) Error() string {
	return e.Source.Loc() + ": " + e.Err.Error()
}

// Wrap wraps err in a ErrorAt or returns err if err is nil, already a
// ErrorAt, or caller info is not available.
//
// XXX: Refactor to wrap.Err(error, ...slog.Attr).  Often want to add attributes for the top level logger.
func Wrap(err error) error {
	// Nothing to do
	if err == nil {
		return nil
	}

	// Already a holos error no need to do anything.
	var errAt *ErrorAt
	if errors.As(err, &errAt) {
		return err
	}

	// Try to wrap err with caller info
	if _, file, line, ok := runtime.Caller(1); ok {
		return &ErrorAt{
			Err: err,
			Source: Source{
				File: file,
				Line: line,
			},
		}
	}

	return err
}

// Log logs err with Source location if Err is a ErrorAt
func Log(log *slog.Logger, ctx context.Context, level slog.Level, err error, msg string, args ...any) {
	var errAt *ErrorAt
	if ok := errors.As(err, &errAt); ok {
		args = append(args,
			slog.String("err", errAt.Unwrap().Error()),
			slog.String("loc", errAt.Source.Loc()),
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

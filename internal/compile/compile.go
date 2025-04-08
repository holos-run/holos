// Package compile compiles BuildPlan resources by reading json encoded data
// from a reader, unmarshaling the data into a Component, building a CUE
// instance injecting the Component as a tag, then exporting a BuildPlan and
// marshalling the result to a writer represented as a stream of json objects.
// Each input component maps to one output json object in the stream.
package compile

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
)

// New returns a new BuildPlan Compiler.
func New(opts ...Option) *Compiler {
	c := &Compiler{
		r: os.Stdin,
		w: os.Stdout,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

type Option func(c *Compiler)

func WithReader(r io.Reader) Option {
	return func(c *Compiler) { c.r = r }
}

func WithWriter(w io.Writer) Option {
	return func(c *Compiler) { c.w = w }
}

type Compiler struct {
	r io.Reader
	w io.Writer
}

func (c *Compiler) Run(ctx context.Context) error {
	dec := json.NewDecoder(c.r)
	enc := json.NewEncoder(c.w)
	slog.DebugContext(ctx, "entering read loop")

	for {
		select {
		case <-ctx.Done():
			return errors.Wrap(ctx.Err())
		default:
		}

		var raw json.RawMessage
		err := dec.Decode(&raw)
		if err == io.EOF {
			slog.DebugContext(ctx, "received: eof", "eof", true)
			return nil
		}

		var meta holos.TypeMeta
		err = json.Unmarshal(raw, &meta)
		if err != nil {
			return errors.Format("could not unmarshal input: %w", err)
		}
		slog.DebugContext(ctx, fmt.Sprintf("received: %+v", meta), "meta", meta)

		err = enc.Encode(raw)
		if err != nil {
			return errors.Format("could not marshal output: %w", err)
		}
	}
}

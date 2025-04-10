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
	"time"

	componentPkg "github.com/holos-run/holos/internal/component"
	"github.com/holos-run/holos/internal/component/v1alpha6"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
)

// New returns a new BuildPlan Compiler.
func New() *Compiler {
	return &Compiler{
		R:        os.Stdin,
		W:        os.Stdout,
		Encoding: "json",
	}
}

// Compiler reads BuildPlanRequest objects from R and exports
// BuildPlanResponse objects to W.
type Compiler struct {
	R io.Reader
	W io.Writer
	// Encoding specifies "json" (default) or "yaml" format output.
	Encoding string
}

// TODO: Define a BuildPlanRequest message which takes a component.  The
// v1alpha6 core.Component fields aren't a good fit for this internal message.
// We need to carry the platform root, WriteTo, and temp directory.  We have the
// path from the component.
func (c *Compiler) Run(ctx context.Context) error {
	epoch := time.Now()
	dec := json.NewDecoder(c.R)
	enc, err := holos.NewSequentialEncoder(c.Encoding, c.W)
	if err != nil {
		return errors.Wrap(err)
	}

	// platform cue module root directory
	root, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err)
	}

	slog.DebugContext(ctx, "entering read loop")
	for idx := 0; ; idx++ {
		log := slog.With("idx", idx)
		select {
		case <-ctx.Done():
			return errors.Wrap(ctx.Err())
		default:
		}

		var raw json.RawMessage
		err := dec.Decode(&raw)
		if err == io.EOF {
			duration := time.Since(epoch)
			msg := fmt.Sprintf("received eof: exiting: total runtime %.3fs", duration.Seconds())
			log.DebugContext(ctx, msg, "eof", true, "seconds", duration.Seconds())
			return nil
		}

		start := time.Now()

		var meta holos.TypeMeta
		err = json.Unmarshal(raw, &meta)
		if err != nil {
			return errors.Format("could not unmarshal type meta: %w", err)
		}
		log.DebugContext(ctx, fmt.Sprintf("received: %+v", meta), "meta", meta)

		writeTo := "deploy"
		// platform component from the components field of the platform resource.
		var pc holos.Component
		switch meta.APIVersion {
		case "v1alpha6":
			var com v1alpha6.Component
			err = json.Unmarshal(raw, &com)
			if err != nil {
				return errors.Format("could not unmarshal component: %w", err)
			}
			if com.WriteTo != "" {
				writeTo = com.WriteTo
			}
			pc = &com
		default:
			log.ErrorContext(ctx, fmt.Sprintf("unsupported api version: %+v (ignored)", meta), "meta", meta)
			continue
		}

		// Produce the build plan.
		component := componentPkg.New(root, pc.Path(), componentPkg.NewConfig())
		tm, err := component.TypeMeta()
		if err != nil {
			return errors.Wrap(err)
		}
		opts := holos.NewBuildOpts(root, pc.Path(), writeTo, "${TMPDIR_PLACEHOLDER}")

		// Component name, label, annotations passed via tags to cue.
		tags, err := pc.Tags()
		if err != nil {
			return errors.Wrap(err)
		}
		opts.Tags = tags

		bp, err := component.BuildPlan(tm, opts)
		if err != nil {
			return errors.Wrap(err)
		}

		if err := bp.Export(idx, enc); err != nil {
			return errors.Format("could not marshal output: %w", err)
		}

		duration := time.Since(start)
		log.DebugContext(ctx, fmt.Sprintf("compile time: %.3fs", duration.Seconds()))
	}
}

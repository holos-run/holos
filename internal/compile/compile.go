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
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
)

// BuildPlanRequest represents the complete context necessary to produce a
// BuildPlan.  BuildPlanRequest is the primary input to the holos compile
// command, read from standard input.  Provided by the holos render platform
// command.
type BuildPlanRequest struct {
	holos.TypeMeta
	Root    string `json:"root,omitempty" yaml:"root,omitempty"`
	Leaf    string `json:"leaf,omitempty" yaml:"leaf,omitempty"`
	WriteTo string `json:"writeTo,omitempty" yaml:"writeTo,omitempty"`
	TempDir string `json:"tempDir,omitempty" yaml:"tempDir,omitempty"`
	Tags    []string
}

type BuildPlanResponse struct {
	holos.TypeMeta
	BuildPlan json.RawMessage `json:"buildPlan,omitempty" yaml:"buildPlan,omitempty"`
}

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

// Run reads a BuildPlanRequest from R and compiles a BuildPlan into a
// BuildPlanResponse on W.  R and W are usually connected to stdin and stdout.
func (c *Compiler) Run(ctx context.Context) error {
	epoch := time.Now()
	dec := json.NewDecoder(c.R)
	enc, err := holos.NewSequentialEncoder(c.Encoding, c.W)
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

		var req BuildPlanRequest

		// platform component from the components field of the platform resource.
		switch meta.APIVersion {
		case "v1alpha6":
			switch meta.Kind {
			case holos.BuildPlanRequest:
				err = json.Unmarshal(raw, &req)
				if err != nil {
					return errors.Format("could not unmarshal %+v: %w", meta, err)
				}
			default:
				return errors.Format("unsupported kind: %+v (ignored)", meta)
			}
		default:
			return errors.Format("unsupported api version: %+v (ignored)", meta)
		}

		// Produce the build plan.
		component := componentPkg.New(req.Root, req.Leaf, componentPkg.NewConfig())
		tm, err := component.TypeMeta()
		if err != nil {
			return errors.Wrap(err)
		}
		// TODO(jjm): Decide how to handle the temp directory.
		opts := holos.NewBuildOpts(req.Root, req.Leaf, req.WriteTo, "${TMPDIR_PLACEHOLDER}")

		// Component name, label, annotations passed via tags to cue.
		opts.Tags = req.Tags

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

// Package compile compiles BuildPlan resources by reading json encoded data
// from a reader, unmarshaling the data into a Component, building a CUE
// instance injecting the Component as a tag, then exporting a BuildPlan and
// marshalling the result to a writer represented as a stream of json objects.
// Each input component maps to one output json object in the stream.
package compile

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"time"

	componentPkg "github.com/holos-run/holos/internal/component"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/logger"
	"golang.org/x/sync/errgroup"
)

// BuildPlanRequest represents the complete context necessary to produce a
// BuildPlan.  BuildPlanRequest is the primary input to the holos compile
// command, read from standard input.  Provided by the holos render platform
// command.
type BuildPlanRequest struct {
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty" yaml:"kind,omitempty"`
	Root       string `json:"root,omitempty" yaml:"root,omitempty"`
	Leaf       string `json:"leaf,omitempty" yaml:"leaf,omitempty"`
	WriteTo    string `json:"writeTo,omitempty" yaml:"writeTo,omitempty"`
	TempDir    string `json:"tempDir,omitempty" yaml:"tempDir,omitempty"`
	Tags       []string
}

type BuildPlanResponse struct {
	APIVersion string          `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string          `json:"kind,omitempty" yaml:"kind,omitempty"`
	RawMessage json.RawMessage `json:"rawMessage,omitempty" yaml:"rawMessage,omitempty"`
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
	decoder := json.NewDecoder(c.R)
	encoder, err := holos.NewSequentialEncoder(c.Encoding, c.W)
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
		err := decoder.Decode(&raw)
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

		if err := bp.Export(idx, encoder); err != nil {
			return errors.Format("could not marshal output: %w", err)
		}

		duration := time.Since(start)
		log.DebugContext(ctx, fmt.Sprintf("compile time: %.3fs", duration.Seconds()))
	}
}

type task struct {
	idx  int
	req  BuildPlanRequest
	resp BuildPlanResponse
}

func Compile(ctx context.Context, concurrency int, reqs []BuildPlanRequest) (resp []BuildPlanResponse, err error) {
	concurrency = min(len(reqs), max(1, concurrency))
	resp = make([]BuildPlanResponse, len(reqs))

	g, ctx := errgroup.WithContext(ctx)
	tasks := make(chan task)

	// Producer
	g.Go(func() error {
		for idx, req := range reqs {
			tsk := task{
				idx:  idx,
				req:  req,
				resp: BuildPlanResponse{},
			}
			select {
			case <-ctx.Done():
				return errors.Wrap(ctx.Err())
			case tasks <- tsk:
				slog.DebugContext(ctx, fmt.Sprintf("producer producing task seq=%d component=%s tags=%+v", tsk.idx, tsk.req.Leaf, tsk.req.Tags))
			}
		}
		slog.DebugContext(ctx, fmt.Sprintf("producer finished: closing tasks channel"))
		close(tasks)
		return nil
	})

	// Consumers
	for id := range concurrency {
		g.Go(func() error {
			return compiler(ctx, id, tasks, resp)
		})
	}

	err = errors.Wrap(g.Wait())
	return
}

func compiler(ctx context.Context, id int, tasks chan task, resp []BuildPlanResponse) error {
	log := logger.FromContext(ctx).With("id", id)

	// Start the sub-process
	exe, err := os.Executable()
	if err != nil {
		return errors.Wrap(err)
	}

	cmd := exec.CommandContext(ctx, exe, "compile")

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return errors.Format("could not attach to stdin for worker %d: %w", id, err)
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		_ = stdinPipe.Close()
		return errors.Format("could not attach to stdout for worker %d: %w", id, err)
	}
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	if err := cmd.Start(); err != nil {
		_ = stdinPipe.Close()
		return errors.Format("could not start worker %d: %w", id, err)
	}
	pid := cmd.Process.Pid
	msg := fmt.Sprintf("compiler id=%d pid=%d", id, pid)
	log.DebugContext(ctx, fmt.Sprintf("%s: started", msg))

	defer func() {
		stdinPipe.Close()
		if err := cmd.Wait(); err != nil {
			log.ErrorContext(ctx, fmt.Sprintf("%s: exited uncleanly: %s", msg, err), "err", err, "stderr", stderrBuf.String())
		}
	}()

	encoder := json.NewEncoder(stdinPipe)
	decoder := json.NewDecoder(stdoutPipe)

	for {
		select {
		case <-ctx.Done():
			return errors.Wrap(ctx.Err())
		case tsk, ok := <-tasks:
			if !ok {
				log.DebugContext(ctx, fmt.Sprintf("%s: tasks channel closed: returning normally", msg))
				return nil
			}
			log.DebugContext(ctx, fmt.Sprintf("%s: encoding request seq=%d", msg, tsk.idx))
			if err := encoder.Encode(tsk.req); err != nil {
				return errors.Format("could not encode request for %s: %w", msg, err)
			}
			log.DebugContext(ctx, fmt.Sprintf("%s: decoding response seq=%d", msg, tsk.idx))
			if err := decoder.Decode(&resp[tsk.idx].RawMessage); err != nil {
				return errors.Format("could not decode response from %s: %w\n%s", msg, err, stderrBuf.String())
			}
			log.DebugContext(ctx, fmt.Sprintf("%s: ok finished task seq=%d", msg, tsk.idx))
		}
	}
}

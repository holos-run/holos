package util

import (
	"bytes"
	"context"
	"os"
	"os/exec"

	"github.com/holos-run/holos/pkg/logger"
)

// runResult holds the stdout and stderr of a command.
type RunResult struct {
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
}

// RunCmd runs a command within a context, captures its output, provides debug
// logging, and returns the result.
// Example:
//
//	result, err := RunCmd(ctx, "echo", "hello")
//	if err != nil {
//	  return wrapper.Wrap(err)
//	}
//	fmt.Println(result.Stdout.String())
//
// Output:
//
//	"hello\n"
func RunCmd(ctx context.Context, name string, args ...string) (result RunResult, err error) {
	result = RunResult{
		Stdout: new(bytes.Buffer),
		Stderr: new(bytes.Buffer),
	}
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = result.Stdout
	cmd.Stderr = result.Stderr
	log := logger.FromContext(ctx)
	log.DebugContext(ctx, "running: "+name, "name", name, "args", args)
	err = cmd.Run()
	return result, err
}

// RunInteractiveCmd runs a command within a context but allows the command to
// accept stdin interactively from the user. The caller is expected to handle
// errors.
func RunInteractiveCmd(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log := logger.FromContext(ctx)
	log.DebugContext(ctx, "running: "+name, "name", name, "args", args)
	return cmd.Run()
}

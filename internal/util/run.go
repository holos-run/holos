package util

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/holos-run/holos/internal/logger"
)

var mu sync.Mutex

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
	command := fmt.Sprintf("%s '%s'", name, strings.Join(args, "' '"))
	log.DebugContext(ctx, "running command: "+command, "name", name, "args", args)
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("could not run command:\n\t%s\n\t%w", command, err)
	}
	return result, err
}

// RunCmdW calls RunCmd and copies the result stderr to w if there is an error.
func RunCmdW(ctx context.Context, w io.Writer, name string, args ...string) (result RunResult, err error) {
	result, err = RunCmd(ctx, name, args...)
	if err != nil {
		mu.Lock()
		defer mu.Unlock()
		_, err2 := io.Copy(w, result.Stderr)
		if err2 != nil {
			err = fmt.Errorf("could not copy stderr: %s: %w", err2.Error(), err)
		}
	}
	return result, err
}

// RunCmdA calls RunCmd and always copies the result stderr to w.
func RunCmdA(ctx context.Context, w io.Writer, name string, args ...string) (result RunResult, err error) {
	result, err = RunCmd(ctx, name, args...)
	mu.Lock()
	defer mu.Unlock()
	if _, err2 := io.Copy(w, result.Stderr); err2 != nil {
		err = fmt.Errorf("could not copy stderr: %s: %w", err2.Error(), err)
	}
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

package cli

import (
	"bytes"
	"github.com/holos-run/holos/pkg/config"
	"github.com/holos-run/holos/pkg/logger"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"testing"
)

func newCommand() (*cobra.Command, *bytes.Buffer) {
	var b bytes.Buffer
	cmd := New(config.New(os.Stdout, &b))
	return cmd, &b
}

func TestNewRoot(t *testing.T) {
	cmd, _ := newCommand()
	if err := cmd.Execute(); err != nil {
		t.Fatalf("could not execute: %v", err)
	}
}

type argsTestCase struct {
	format string
	level  string
	drop   string
}

func (tc argsTestCase) args() []string {
	return []string{"--log-format", tc.format, "--log-level", tc.level}
}

func TestValidArgs(t *testing.T) {
	t.Parallel()
	formats := []string{"text", "json"}
	levels := []string{"debug", "info", "warn", "error"}
	drops := []string{"version"}
	for _, format := range formats {
		for _, level := range levels {
			for _, drop := range drops {
				tc := argsTestCase{format, level, drop}
				t.Run(strings.Join(tc.args(), " "), func(t *testing.T) {
					cmd, _ := newCommand()
					cmd.SetArgs(tc.args())
					if err := cmd.Execute(); err != nil {
						t.Errorf("unexpected error: %v", err)
					}
				})
			}
		}
	}
}

func TestLogOutput(t *testing.T) {
	// Lifecycle message is always displayed
	cmd, b := newCommand()
	cmd.SetArgs([]string{"--log-level=debug"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("could not execute: %v", err)
	}
	stderr := b.String()
	if !strings.Contains(stderr, "config lifecycle") {
		t.Fatalf("lifecycle message missing: stderr: %v", stderr)
	}
}

func TestLogDrop(t *testing.T) {
	// Log attributes can be filtered out by the user
	cmd, b := newCommand()
	cmd.SetArgs([]string{"--log-level=debug", "--log-drop=version,another"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("could not execute: %v", err)
	}
	stderr := b.String()
	if strings.Contains(stderr, "version") {
		t.Fatalf("want version dropped got: %v", stderr)
	}
}

func TestInvalidArgs(t *testing.T) {
	invalidArgs := [][]string{
		{"--log-format=yaml"},
		{"--log-level=something"},
	}
	for _, args := range invalidArgs {
		var b bytes.Buffer
		cmd := New(config.New(os.Stdout, &b))
		cmd.SetArgs(args)
		err := cmd.Execute()
		if err == nil {
			t.Fatalf("expected error from args: %v", args)
		}
	}
}

func TestLoggerFromContext(t *testing.T) {
	cmd, b := newCommand()
	if err := cmd.Execute(); err != nil {
		t.Fatalf("could not execute: %v", err)
	}
	log := logger.FromContext(cmd.Context())
	want := "612c48a3-30c9-44e3-b8b6-394191d99935"
	log.Info(want)
	have := b.String()
	if !strings.Contains(have, want) {
		t.Fatalf("want: %v have: %v", want, have)
	}

}

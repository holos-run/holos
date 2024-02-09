package cli

import (
	"bytes"
	"github.com/holos-run/holos/pkg/config"
	"github.com/holos-run/holos/pkg/logger"
	"github.com/holos-run/holos/pkg/version"
	"github.com/spf13/cobra"
	"strings"
	"testing"
)

func newCommand() (*cobra.Command, *bytes.Buffer) {
	var b1, b2 bytes.Buffer
	// discard stdout for now, it's a bunch of usage messages.
	cmd := New(config.New(config.Stdout(&b1), config.Stderr(&b2)))
	return cmd, &b2
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
	have := strings.TrimSpace(b.String())
	want := "finalized config from flags"
	if !strings.Contains(have, want) {
		t.Fatalf("have does not contain want\n\thave: %#v\n\twant: %#v", have, want)
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
		cmd := New(config.New(config.Stdout(&b)))
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

func TestVersion(t *testing.T) {
	var b bytes.Buffer
	cmd := New(config.New(config.Stdout(&b)))
	cmd.SetOut(&b)
	cmd.SetArgs([]string{"--version"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("could not execute: %v", err)
	}
	want := version.Version + "\n"
	have := b.String()
	if want != have {
		t.Fatalf("want: %v have: %v", want, have)
	}
}

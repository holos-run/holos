package config

import (
	"bytes"
	"os"
	"testing"
)

func newConfig() (*Config, *bytes.Buffer) {
	var b bytes.Buffer
	c := New(os.Stdout, &b)
	return c, &b
}

func TestConfigFinalize(t *testing.T) {
	cfg, _ := newConfig()
	if err := cfg.Finalize(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.finalized {
		t.Fatalf("not finalized")
	}
}

func TestConfigFinalizeTwice(t *testing.T) {
	cfg, _ := newConfig()
	if err := cfg.Finalize(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := cfg.Finalize(); err == nil {
		t.Fatalf("want error got nil")
	}
}

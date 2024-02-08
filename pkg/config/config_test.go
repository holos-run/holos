package config

import (
	"bytes"
	"testing"
)

// newConfig returns a new *Config with stderr wired to a bytes.Buffer.
func newConfig() (cfg *Config, stderr *bytes.Buffer) {
	var b bytes.Buffer
	return New(Stderr(&b)), &b
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
	cfg, stderr := newConfig()
	if err := cfg.Finalize(); err != nil {
		t.Fatalf("want: %#v have: %#v", nil, err)
	}
	if err := cfg.Finalize(); err == nil {
		t.Fatalf("want: error have: %#v", err)
	} else {
		want := "could not finalize: already finalized"
		have := err.Error()
		if want != have {
			t.Fatalf("want: %#v have: %#v", want, have)
		}
	}
	want := ""
	have := stderr.String()
	if want != have {
		t.Fatalf("want: %#v have: %#v", want, have)
	}
}

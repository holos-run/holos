package errors_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/holos-run/holos/pkg/errors"
)

const prefix = "errors_test.go:"

type record struct {
	Source errors.Source `json:"source"`
	Err    string        `json:"err"`
	Loc    string        `json:"loc"`
	Msg    string        `json:"msg"`
}

func TestLog(t *testing.T) {
	testCases := []error{
		errors.Wrap(fmt.Errorf("when wrapped")),
		fmt.Errorf("when not wrapped"),
	}

	for _, err := range testCases {
		t.Run(err.Error(), func(t *testing.T) {
			var b bytes.Buffer
			log := slog.New(slog.NewJSONHandler(&b, &slog.HandlerOptions{AddSource: true}))
			errors.Log(log, context.Background(), slog.LevelError, err, "handled")
			var rec record
			if err := json.Unmarshal(b.Bytes(), &rec); err != nil {
				t.Fatalf("want json have error: %v", err)
			}

			have := rec.Loc
			want := prefix
			if !strings.HasPrefix(have, want) {
				t.Fatalf("missing prefix:\n\thave: (%v)\n\twant: (%v)", have, want)
			}

			have = rec.Msg
			want = "handled"
			if !strings.HasSuffix(have, want) {
				t.Fatalf("missing suffix:\n\thave: (%v)\n\twant: (%v)", have, want)
			}
		})
	}
}

func TestLogNotWrapped(t *testing.T) {
	want := prefix
	err := fmt.Errorf("boom")
	var b bytes.Buffer
	log := slog.New(slog.NewJSONHandler(&b, &slog.HandlerOptions{AddSource: true}))
	errors.Log(log, context.Background(), slog.LevelError, err, "handled")
	var rec record
	if err := json.Unmarshal(b.Bytes(), &rec); err != nil {
		t.Fatalf("unexpected error:\n\thave: (%v)\n\twant: (%v)", err, nil)
	}

	have := rec.Loc
	if !strings.HasPrefix(have, want) {
		t.Fatalf("missing prefix:\n\thave: (%v)\n\twant: (%v)", have, want)
	}

	have = rec.Msg
	want = "handled"
	if !strings.HasSuffix(have, want) {
		t.Fatalf("missing suffix:\n\thave: (%v)\n\twant: (%v)", have, want)
	}
}

func TestUnwrap(t *testing.T) {
	err := errors.Wrap(fmt.Errorf("boom"))
	want := prefix
	have := err.Error()

	if !strings.HasPrefix(have, want) {
		t.Fatalf("missing prefix:\n\thave: (%v)\n\twant: (%v)", have, want)
	}

	want = "boom"
	if !strings.HasSuffix(have, want) {
		t.Fatalf("missing suffix:\n\thave: (%v)\n\twant: (%v)", have, want)
	}
}

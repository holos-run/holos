package wrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"testing"
)

type record struct {
	Source Source `json:"source"`
	Err    string `json:"err"`
	Loc    string `json:"loc"`
	Msg    string `json:"msg"`
}

func TestLog(t *testing.T) {
	testCases := []error{
		Wrap(fmt.Errorf("when wrapped")),
		fmt.Errorf("when not wrapped"),
	}

	for _, err := range testCases {
		t.Run(err.Error(), func(t *testing.T) {
			var b bytes.Buffer
			log := slog.New(slog.NewJSONHandler(&b, &slog.HandlerOptions{AddSource: true}))
			Log(log, context.Background(), slog.LevelError, err, "handled")
			var rec record
			if err := json.Unmarshal(b.Bytes(), &rec); err != nil {
				t.Fatalf("want json have error: %v", err)
			}

			have := rec.Loc
			want := "wrap_test.go:"
			if !strings.HasPrefix(have, want) {
				t.Fatalf("want: %v have: %v", want, have)
			}

			have = rec.Msg
			want = "handled"
			if !strings.HasSuffix(have, want) {
				t.Fatalf("want: %v have: %v", want, have)
			}
		})
	}
}

func TestLogNotWrapped(t *testing.T) {
	want := "wrap_test.go:"
	err := fmt.Errorf("boom")
	var b bytes.Buffer
	log := slog.New(slog.NewJSONHandler(&b, &slog.HandlerOptions{AddSource: true}))
	Log(log, context.Background(), slog.LevelError, err, "handled")
	var rec record
	if err := json.Unmarshal(b.Bytes(), &rec); err != nil {
		t.Fatalf("want json have error: %v", err)
	}

	have := rec.Loc
	if !strings.HasPrefix(have, want) {
		t.Fatalf("want: %v have: %v", want, have)
	}

	have = rec.Msg
	want = "handled"
	if !strings.HasSuffix(have, want) {
		t.Fatalf("want: %v have: %v", want, have)
	}
}

func TestUnwrap(t *testing.T) {
	err := Wrap(fmt.Errorf("boom"))
	want := "wrap_test.go:"
	have := err.Error()

	if !strings.HasPrefix(have, want) {
		t.Fatalf("want: %v have: %v", want, have)
	}

	want = "boom"
	if !strings.HasSuffix(have, want) {
		t.Fatalf("want: %v have: %v", want, have)
	}
}

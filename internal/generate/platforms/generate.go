package platforms

//go:generate rm -rf cue.mod/gen/github.com/holos-run/holos/api/v1alpha1
//go:generate cue get go github.com/holos-run/holos/api/v1alpha1/...

//go generate rm -rf cue.mod/gen/github.com/holos-run/holos/api/core
//go:generate cue get go github.com/holos-run/holos/api/core/...

//go generate rm -rf cue.mod/gen/github.com/holos-run/holos/api/meta
//go:generate cue get go github.com/holos-run/holos/api/meta/...

//go:generate touch ../platform.go

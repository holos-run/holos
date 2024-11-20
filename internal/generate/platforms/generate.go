package platforms

// TODO: Remove env GODEBUG=gotypesalias=0 when cue 0.11 is released and used.
// See: https://github.com/cue-lang/cue/issues/3539

//go generate rm -rf cue.mod/gen/github.com/holos-run/holos/api/core
//go:generate env GODEBUG=gotypesalias=0 cue get go github.com/holos-run/holos/api/core/...

//go generate rm -rf cue.mod/gen/github.com/holos-run/holos/api/meta
//go:generate env GODEBUG=gotypesalias=0 cue get go github.com/holos-run/holos/api/meta/...

//go generate rm -rf cue.mod/gen/github.com/holos-run/holos/api/author
//go:generate env GODEBUG=gotypesalias=0 cue get go github.com/holos-run/holos/api/author/...

//go generate rm -rf cue.mod/gen/github.com/holos-run/holos/service/gen/holos/object
//go:generate cue import ../../../service/holos/object/v1alpha1/object.proto -o cue.mod/gen/github.com/holos-run/holos/service/gen/holos/object/v1alpha1/object.proto_gen.cue -I ../../../proto -f
//go:generate rm -f cue.mod/gen/github.com/holos-run/holos/service/gen/holos/object/v1alpha1/object.pb_go_gen.cue

//go:generate touch ../platform.go

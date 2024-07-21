package platforms

//go:generate rm -rf cue.mod/gen/github.com/holos-run/holos/api/v1alpha1
//go:generate cue get go github.com/holos-run/holos/api/v1alpha1/...

//go generate rm -rf cue.mod/gen/github.com/holos-run/holos/api/core
//go:generate cue get go github.com/holos-run/holos/api/core/...

//go generate rm -rf cue.mod/gen/github.com/holos-run/holos/api/meta
//go:generate cue get go github.com/holos-run/holos/api/meta/...

//go generate rm -rf cue.mod/gen/github.com/holos-run/holos/service/gen/holos/object
//go:generate cue import ../../../service/holos/object/v1alpha1/object.proto -o cue.mod/gen/github.com/holos-run/holos/service/gen/holos/object/v1alpha1/object.proto_gen.cue -I ../../../proto -f
//go:generate rm -f cue.mod/gen/github.com/holos-run/holos/service/gen/holos/object/v1alpha1/object.pb_go_gen.cue

//go:generate touch ../platform.go

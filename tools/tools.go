//go:build tools

package tools

// Refer to "How can I track tool dependencies for a module?"
// https://go.dev/wiki/Modules

import (
	_ "connectrpc.com/connect/cmd/protoc-gen-connect-go"
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/fullstorydev/grpcurl/cmd/grpcurl"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
	_ "honnef.co/go/tools/cmd/staticcheck"
)

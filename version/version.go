package version

import (
	_ "embed"
	"fmt"
	"runtime"
	"strings"
)

//go:embed embedded/major
var Major string

//go:embed embedded/minor
var Minor string

//go:embed embedded/patch
var Patch string

// Version of this module for example "0.70.0"
var Version = strings.ReplaceAll(Major+"."+Minor+"."+Patch, "\n", "")

// GitDescribe is intended as the output of git describe --tags plus DIRTY if dirty.
var GitDescribe string

// GitCommit defined dynamically by the Makefile
var GitCommit string
var GitTreeState string

// BuildDate as an iso-8601 string with seconds precision.
var BuildDate = ""

// GoVersion returns the version of the go runtime used to compile the binary
var GoVersion = runtime.Version()

// OsArch returns the os and arch used to build the binary
var OsArch = fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)

type VersionInfo struct {
	Version      string `yaml:"Version"`
	GitCommit    string `yaml:"Git Commit"`
	GitTreeState string `yaml:"Git Tree State"`
	GoVersion    string `yaml:"Go Version"`
	BuildDate    string `yaml:"Build Date"`
	Os           string `yaml:"Os"`
	Arch         string `yaml:"Arch"`
}

type VersionBrief struct {
	Version      string `json:"version"`
	GitCommit    string `json:"commit"`
	GitTreeState string `json:"tree"`
}

// NewVersionInfo returns a filled in struct for marshaling
func NewVersionInfo() VersionInfo {
	return VersionInfo{
		Version:      Version,
		GitCommit:    GitCommit,
		GitTreeState: GitTreeState,
		GoVersion:    GoVersion,
		BuildDate:    BuildDate,
		Os:           runtime.GOOS,
		Arch:         runtime.GOARCH,
	}
}

func NewVersionBrief() VersionBrief {
	return VersionBrief{
		Version:      Version,
		GitCommit:    GitCommit,
		GitTreeState: GitTreeState,
	}
}

func GetVersion() string {
	if GitDescribe != "" {
		return strings.Replace(GitDescribe, "v", "", 1)
	}
	return Version
}

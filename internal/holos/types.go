package holos

import (
	"fmt"
	"os"
	"strings"
)

// StringSlice represents zero or more flag values.
type StringSlice []string

// String implements the flag.Value interface.
func (i *StringSlice) String() string {
	return fmt.Sprint(*i)
}

// Type implements the pflag.Value interface and describes the type.
func (i *StringSlice) Type() string {
	return "strings"
}

// Set implements the flag.Value interface.
func (i *StringSlice) Set(value string) error {
	for _, str := range strings.Split(value, ",") {
		*i = append(*i, str)
	}
	return nil
}

type feature string

const BuildFeature = feature("BUILD")
const ServerFeature = feature("SERVER")
const ClientFeature = feature("CLIENT")
const PreflightFeature = feature("PREFLIGHT")
const GenerateComponentFeature = feature("GENERATE_COMPONENT")
const SecretsFeature = feature("SECRETS")

// Flagger is the interface to check if an experimental feature is enabled.
type Flagger interface {
	Flag(name feature) bool
}

type EnvFlagger struct{}

// Flag returns true if feature name is enabled.
func (e *EnvFlagger) Flag(name feature) bool {
	return os.Getenv(fmt.Sprintf("HOLOS_FEATURE_%s", name)) != ""
}

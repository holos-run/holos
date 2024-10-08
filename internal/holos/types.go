package holos

import (
	"fmt"
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

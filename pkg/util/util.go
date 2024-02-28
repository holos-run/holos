package util

// EnsureNewline adds a trailing newline if not already there.
func EnsureNewline(b []byte) []byte {
	if len(b) > 0 && b[len(b)-1] != '\n' {
		b = append(b, '\n')
	}
	return b
}

package types

#Target: {
	// Ensure name has a concrete value so json.Marshal works.  See
	// https://github.com/holos-run/holos/issues/348
	name: string | *""
}

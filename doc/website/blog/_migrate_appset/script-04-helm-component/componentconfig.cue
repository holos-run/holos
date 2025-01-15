package holos

#ComponentConfig: {
	// Inject the output base directory from platform component parameters.
	OutputBaseDir: string | *"" @tag(outputBaseDir, type=string)
}

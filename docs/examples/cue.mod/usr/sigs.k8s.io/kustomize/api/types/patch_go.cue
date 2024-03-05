package types

#Patch: {
	// Path is a relative file path to the patch file.
	path?: string @go(Path)

	// Patch is the content of a patch.
	patch?: string @go(Patch)

	// Target points to the resources that the patch is applied to
	target?: #Target | #Selector @go(Target,*Selector)

	// Options is a list of options for the patch
	options?: {[string]: bool} @go(Options,map[string]bool)
}

package holos

// Override these values with cue build tags.  See organization-jeff.cue for an
// example.
Organization: #Organization & {
	DisplayName: string | *"Kargo Demo"
	Domain:      string | *"holos.localhost"
	RepoURL:     string | *"https://github.com/holos-run/kargo-demo.git"
}

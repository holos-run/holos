package holos

// Istio meshconfig
// TODO: Generate per-project extauthz providers.
_MeshConfig: (#MeshConfig & {projects: _Projects}).config

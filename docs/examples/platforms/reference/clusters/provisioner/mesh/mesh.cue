package holos

// Components under this directory are part of this collection
#InputKeys: project: "mesh"

// Shared dependencies for all components in this collection.
#DependsOn: _Namespaces

// Common Dependencies
_CertManager: CertManager: name: "\(#InstancePrefix)-certmanager"
_Namespaces: Namespaces: name:   "\(#StageName)-secrets-namespaces"

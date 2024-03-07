package holos

// Components under this directory are part of this collection
#InputKeys: project: "secrets"

// Shared dependencies for all components in this collection.
#DependsOn: _Namespaces

// Common Dependencies
_Namespaces: Namespaces: name: "\(#StageName)-secrets-namespaces"
_ESO: ESO: name:               "\(#InstancePrefix)-eso"
_ESOCreds: ESOCreds: name:     "\(#InstancePrefix)-eso-creds-refresher"

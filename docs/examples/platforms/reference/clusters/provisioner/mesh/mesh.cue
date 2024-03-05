package holos

// Components under this directory are part of this collection
#InputKeys: project: "mesh"

// Shared dependencies for all components in this collection.
#DependsOn: _Namespaces

// Common Dependencies
_Namespaces: Namespaces: name:     "\(#StageName)-secrets-namespaces"
_CertManager: CertManager: name:   "\(#InstancePrefix)-certmanager"
_LetsEncrypt: LetsEncrypt: name:   "\(#InstancePrefix)-letsencrypt"
_Certificates: Certificates: name: "\(#InstancePrefix)-certificates"

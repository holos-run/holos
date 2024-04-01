package holos

// Ingress Gateway default auth proxy
#MeshConfig: extensionProviderMap: ingressauth: envoyExtAuthzHttp: service: #IngressAuthProxy.service

// Istio meshconfig
_MeshConfig: (#MeshConfig & {projects: _Projects}).config

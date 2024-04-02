package holos

// Ingress Gateway default auth proxy
let Provider = _IngressAuthProxy.AuthProxySpec.provider
let Service = _IngressAuthProxy.service
#MeshConfig: extensionProviderMap: (Provider): envoyExtAuthzHttp: service: Service

// Istio meshconfig
_MeshConfig: (#MeshConfig & {projects: _Projects}).config

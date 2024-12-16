package holos

holos: Component.BuildPlan
Component: #Kustomize & {
	Name: "gateway-api"

	KustomizeConfig: Files: {
		"standard/gateway.networking.k8s.io_gatewayclasses.yaml":  _
		"standard/gateway.networking.k8s.io_gateways.yaml":        _
		"standard/gateway.networking.k8s.io_grpcroutes.yaml":      _
		"standard/gateway.networking.k8s.io_httproutes.yaml":      _
		"standard/gateway.networking.k8s.io_referencegrants.yaml": _
	}
}

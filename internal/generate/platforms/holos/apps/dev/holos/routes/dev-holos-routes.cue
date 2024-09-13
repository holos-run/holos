package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).BuildPlan

_AppInfo: spec: component: "routes"

let Objects = {
	Name:      _AppInfo.status.component
	Namespace: #IstioGatewaysNamespace

	Resources: [_]: [_]: metadata: namespace: Namespace

	// HTTPRoute resources
	Resources: (#AppRoute & {AppInfo: _AppInfo}).Resources
}

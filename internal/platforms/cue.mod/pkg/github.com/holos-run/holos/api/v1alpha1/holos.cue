package v1alpha1

// #HolosComponent defines struct fields common to all holos platform
// component types.
#HolosComponent: {
	metadata: name: string
	_NameLengthConstraint: len(metadata.name) & >=1
	...
}

// #KubernetesObjects is a Holos Component BuildPlan which has k8s api objects
// embedded into the build plan itself.
#KubernetesObjects: #HolosComponent & {
	apiVersion: #APIVersion
	kind:       #KubernetesObjectsKind
}

package holos

#InputKeys: component: "cni"
#TargetNamespace: "kube-system"

#HelmChart & {
	namespace: #TargetNamespace
	chart: name: "cni"
	values: #IstioValues
}

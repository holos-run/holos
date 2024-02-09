package holos

// Output schema
{} & #KubernetesObjects & {
	ksObjects: [#Kustomization]
}

#InputKeys: {
	project: "secrets"
	service: "eso"
}

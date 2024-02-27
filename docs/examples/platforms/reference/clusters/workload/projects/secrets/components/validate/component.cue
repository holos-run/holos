package holos

// Validate ESO by syncing a secret with a SecretStore.

#TargetNamespace: "holos-system"

#InputKeys: {
	project:   "secrets"
	component: "validate"
}

#Kustomization: spec: dependsOn: [{name: #InstancePrefix + "-eso"}]

objects: [
	#SecretStore,
	#ExternalSecret & {
		_name: "validate"
		metadata: namespace: #TargetNamespace
		spec: dataFrom: [{extract: key: "ns/" + #TargetNamespace + "/test"}]
	},
]

{} & #KubernetesObjects

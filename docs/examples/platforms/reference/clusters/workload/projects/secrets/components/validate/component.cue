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
		spec: data: [{remoteRef: key: _name}]
	},
]

{} & #KubernetesObjects

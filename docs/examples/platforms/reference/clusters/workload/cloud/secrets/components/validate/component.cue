package holos

// Validate ESO by syncing a secret with a SecretStore.

#TargetNamespace: "holos-system"

#InputKeys: {
	project:   "secrets"
	component: "validate"
}

#Kustomization: spec: dependsOn: [{name: #InstancePrefix + "-eso"}]

#KubernetesObjects & {
	apiObjects: {
		ExternalSecret: validate: #ExternalSecret & {
			_name: "validate"
		}
	}
}

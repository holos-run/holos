package holos

// Component fields are unified with podinfo.cue
Component: {
	// Concrete values are defined in podinfo.cue
	Name:      string
	Namespace: string

	// Resources represents mix-in resources organized as a struct.
	Resources: ExternalSecret: (Name): {
		// Name is consistent with the component name.
		metadata: name: Name
		// Namespace is consistent with the component namespace.
		metadata: namespace: Namespace
		spec: {
			// Ensure the target secret name is consistent.
			target: name: metadata.name
			// Ensure the name in the SecretStore is consistent.
			dataFrom: [{extract: {key: metadata.name}}]
			refreshInterval: "30s"
			secretStoreRef: kind: "SecretStore"
			secretStoreRef: name: "default"
		}
	}
}

package v1alpha5

// Don't allow password values, use fromEnv instead.
#Auth: password: value?: ""

#Transformer: {
	kind: _

	if kind == "Kustomize" {
		kustomize: _
	}

	if kind == "Join" {
		join: _
	}
}

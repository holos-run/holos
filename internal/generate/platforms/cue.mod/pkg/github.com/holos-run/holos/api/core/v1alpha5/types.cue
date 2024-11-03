package v1alpha5

#Transformer: {
	kind: _

	if kind == "Kustomize" {
		kustomize: _
	}

	if kind == "Join" {
		join: _
	}
}

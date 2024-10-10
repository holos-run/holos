package v1alpha4

#Transformer: {
	kind: _

	if kind == "Kustomize" {
		kustomize: _
	}

	if kind == "Join" {
		join: _
	}
}

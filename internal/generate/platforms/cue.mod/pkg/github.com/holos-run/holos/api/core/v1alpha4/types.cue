package v1alpha4

#Component: tags?: [_]: !~"^holos_"

#Transformer: {
	kind: _

	if kind == "Kustomize" {
		kustomize: _
	}

	if kind == "Join" {
		join: _
	}
}

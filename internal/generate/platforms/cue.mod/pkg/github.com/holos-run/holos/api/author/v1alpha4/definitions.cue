package v1alpha4

#Platform: {
	Name:       string | *"no-platform-name"
	Components: _
	Resource: {
		metadata: name: Name
		spec: components: [for x in Components {x}]
	}
}

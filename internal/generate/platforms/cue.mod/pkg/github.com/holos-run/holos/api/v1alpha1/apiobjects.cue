package v1alpha1

// #APIObjects defines the output format for kubernetes api objects.  The holos
// cli expects the yaml representation of each api object in the apiObjectMap
// field.
#APIObjects: {
	// apiObjects represents the un-marshalled form of each kubernetes api object
	// managed by a holos component.
	apiObjects: [Kind=string]: [string]: kind: Kind
	// apiObjectMap holds the marshalled representation of apiObjects
	apiObjectsMap: [string]: [string]: string
}

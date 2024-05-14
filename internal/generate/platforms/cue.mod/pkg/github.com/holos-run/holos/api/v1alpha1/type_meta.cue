package v1alpha1

// #TypeMeta is similar to kubernetes api TypeMeta, but for holos api
// objects such as the Platform config resource.
#TypeMeta: {
	kind:       string @go(Kind)
	apiVersion: string @go(APIVersion)
}

// #ObjectMeta is similar to kubernetes api ObjectMeta, but for holos api
// objects.
#ObjectMeta: {
	name: string @go(Name)
	labels: {[string]: string} @go(Labels,map[string]string)
	annotations: {[string]: string} @go(Annotations,map[string]string)
}

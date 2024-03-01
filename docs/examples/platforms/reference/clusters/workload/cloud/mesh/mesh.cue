package holos

// All components are share this collection
#InputKeys: project: "mesh"

// Shared dependencies for all components in this collection.
#Kustomization: spec: {
	dependsOn: [{name: "\(#StageName)-secrets-namespaces"}, ...]
	targetNamespace: #TargetNamespace
}

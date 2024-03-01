package holos

// Components under this directory are part of this collection
#InputKeys: project: "mesh"

// Shared dependencies for all components in this collection.
#Kustomization: spec: {
	dependsOn: [{name: "\(#StageName)-secrets-namespaces"}, ...]
	targetNamespace: #TargetNamespace
}

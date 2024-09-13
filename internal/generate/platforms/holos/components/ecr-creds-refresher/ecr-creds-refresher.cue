package holos

let NAME = "ecr-creds-refresher"

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).BuildPlan

let AWS_ACCOUNT = _Platform.Model.aws.accountNumber

let Objects = {
	Name:      NAME
	Namespace: "default"

	for Namespace in _Namespaces {
		Resources: ExternalSecret: "\(Namespace.metadata.name)/ecr-creds-\(AWS_ACCOUNT)": #ExternalSecret & {
			metadata: name:      "ecr-creds-\(AWS_ACCOUNT)"
			metadata: namespace: Namespace.metadata.name
		}
	}
}

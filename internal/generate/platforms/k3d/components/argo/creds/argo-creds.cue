package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

let Objects = {
	Name:      "argo-creds"
	Namespace: _ArgoCD.metadata.namespace

	Resources: [_]: [_]: metadata: namespace: Namespace

	Resources: {
		// ssh-keygen -t ed25519 -f sshPrivateKey -m pem -C argocd -N ''
		// echo echo git@github.com:myorg/holos-infra.git > url
		// holos create secret -n argocd --append-hash=false creds-holos-infra --from-file .
		ExternalSecret: "creds-holos-infra": #ExternalSecret & {
			// Labels and annotations are copied over
			metadata: labels: "argocd.argoproj.io/secret-type": "repo-creds"
		}
	}
}

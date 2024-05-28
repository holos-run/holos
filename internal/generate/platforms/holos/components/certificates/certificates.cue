package holos

let Objects = {
	Name:      "certificates"
	Namespace: "istio-ingress"

	Resources: Certificate: _Certificates
}

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

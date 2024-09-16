package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).BuildPlan

let BankName = #BankOfHolos.Name

let CommonLabels = {
	application: BankName
	environment: "development"
	tier:        "backend"
}

let Objects = {
	Name:      "bank-backend-config"
	Namespace: #BankOfHolos.Backend.Namespace

	// Ensure resources go in the correct namespace
	Resources: [_]: [_]: metadata: namespace: Namespace
	Resources: [_]: [_]: metadata: labels:    CommonLabels

	// https://github.com/GoogleCloudPlatform/bank-of-anthos/blob/release/v0.6.5/kubernetes-manifests/userservice.yaml
	Resources: {
		// Allow HTTPRoutes in the ingress gateway namespace to reference Services
		// in this namespace.
		ReferenceGrant: grant: #ReferenceGrant & {
			metadata: namespace: Namespace
		}

		// Include shared resources
		#BankOfHolos.Resources
	}
}

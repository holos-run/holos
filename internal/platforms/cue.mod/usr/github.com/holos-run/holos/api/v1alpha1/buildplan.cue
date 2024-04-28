package v1alpha1

// #BuildPlan is the API contract between CUE and the Holos cli.
// Holos requires CUE to evaluate and provide a valid #BuildPlan.
#BuildPlan: {
	kind:       #BuildPlanKind
	apiVersion: #APIVersion
}

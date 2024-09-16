package holos

import (
	"encoding/yaml"
	ks "sigs.k8s.io/kustomize/api/types"
)

(#Kubernetes & {Name: "argocd-crds"}).BuildPlan

// Holos stages BuildPlan resources as an intermediate step of the rendering
// pipeline.  The purpose is to provide the resources to kustomize for
// post-processing.
let BuildPlanResources = "build-plan-resources.yaml"

let Kustomization = ks.#Kustomization & {
	apiVersion: "kustomize.config.k8s.io/v1beta1"
	kind:       "Kustomization"
	resources: [
		// Kustomize the intermediate build plan resources.
		BuildPlanResources,
		// Mix-in external resources.
		"https://raw.githubusercontent.com/argoproj/argo-cd/v\(#ArgoCD.Version)/manifests/crds/application-crd.yaml",
		"https://raw.githubusercontent.com/argoproj/argo-cd/v\(#ArgoCD.Version)/manifests/crds/applicationset-crd.yaml",
		"https://raw.githubusercontent.com/argoproj/argo-cd/v\(#ArgoCD.Version)/manifests/crds/appproject-crd.yaml",
		// This method also works, but takes about 5 seconds longer each build.
		// "https://github.com/argoproj/argo-cd//manifests/crds/?ref=v\(#ArgoCD.Version)",
	]
}

// Generate a kustomization.yaml directly from CUE so we can provide the correct
// version.
spec: components: kubernetesObjectsList: [{
	// intermediate build plan resources to kustomize.  Necessary to activate the
	// kustomization post-rendering step in holos.
	kustomize: resourcesFile: BuildPlanResources
	kustomize: kustomizeFiles: "kustomization.yaml": yaml.Marshal(Kustomization)
}]

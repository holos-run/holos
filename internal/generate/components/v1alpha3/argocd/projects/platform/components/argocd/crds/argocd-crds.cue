package holos

import (
	"encoding/yaml"
	ks "sigs.k8s.io/kustomize/api/types"
)

(#Kubernetes & {Name: "argocd-crds"}).BuildPlan

let Kustomization = ks.#Kustomization & {
	apiVersion: "kustomize.config.k8s.io/v1beta1"
	kind:       "Kustomization"
	resources: ["https://github.com/argoproj/argo-cd//manifests/crds/?ref=v\(#ArgoCD.Version)"]
}

// Generate a kustomization.yaml directly from CUE so that we can manage the
// correct version of the custom resource definitions.
spec: components: kubernetesObjectsList: [{
	kustomize: kustomizeFiles: "kustomization.yaml": yaml.Marshal(Kustomization)
}]

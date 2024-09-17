package holos

import (
	"encoding/yaml"
	ks "sigs.k8s.io/kustomize/api/types"
)

// Holos stages BuildPlan resources as an intermediate step of the rendering
// pipeline.  The purpose is to provide the resources to kustomize for
// post-processing.
let BuildPlanOutputManifest = "build-plan-output-manifest.yaml"

// Patch istio so it's not constantly out of sync in ArgoCD
let Kustomization = ks.#Kustomization & {
	apiVersion: "kustomize.config.k8s.io/v1beta1"
	kind:       "Kustomization"
	// Kustomize the build plan output.
	resources: [BuildPlanOutputManifest]
	// Patch the the build plan output.
	patches: [for x in KustomizePatches {x}]
}

#KustomizePatches: [ArbitraryLabel=string]: ks.#Patch
let KustomizePatches = #KustomizePatches & {
	validator: {
		target: {
			group:   "admissionregistration.k8s.io"
			version: "v1"
			kind:    "ValidatingWebhookConfiguration"
			name:    "istiod-default-validator"
		}
		let Patch = [{
			op:    "replace"
			path:  "/webhooks/0/failurePolicy"
			value: "Fail"
		}]
		patch: yaml.Marshal(Patch)
	}
}

// Generate a kustomization.yaml directly from CUE so we can provide the correct
// version.
spec: components: helmChartList: [{
	// intermediate build plan resources to kustomize.  Necessary to activate the
	// kustomization post-rendering step in holos.
	kustomize: resourcesFile: BuildPlanOutputManifest
	kustomize: kustomizeFiles: "kustomization.yaml": yaml.Marshal(Kustomization)
}]

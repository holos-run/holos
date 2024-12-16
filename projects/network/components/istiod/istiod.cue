package holos

import (
	"encoding/yaml"
	ks "sigs.k8s.io/kustomize/api/types"
)

// Produce a helm chart build plan.
holos: Component.BuildPlan

Component: #Helm & {
	Name:      "istiod"
	Namespace: Istio.System.Namespace

	Chart: {
		version: Istio.Version
		repository: {
			name: "istio"
			url:  "https://istio-release.storage.googleapis.com/charts"
		}
	}

	KustomizeConfig: Kustomization: patches: [for x in KustomizePatches {x}]

	Values: Istio.Values
}

#KustomizePatches: [ArbitraryLabel=string]: ks.#Patch
let KustomizePatches = #KustomizePatches & {
	validator: {
		target: {
			group:   "admissionregistration.k8s.io"
			version: "v1"
			kind:    "ValidatingWebhookConfiguration"
			name:    "istio-validator-istio-system"
		}
		let Patch = [{
			op:    "replace"
			path:  "/webhooks/0/failurePolicy"
			value: "Fail"
		}]
		patch: yaml.Marshal(Patch)
	}
}

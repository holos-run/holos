@if(!NoKargo)
package holos

import ks "sigs.k8s.io/kustomize/api/types"

// Configure the ArgoCD Application to allow updates from Kargo.  Configure a
// stub kustomization.yaml artifact in the output directory for Kargo to edit.
Component: {
	Name:          _
	OutputBaseDir: _
	_OutPath:      "\(OutputBaseDir)/components/\(Name)"

	_ArgoApplication: {
		metadata: annotations: "kargo.akuity.io/authorized-stage": "\(ProjectName):\(Name)"
		spec: source: {
			path:           "./"
			targetRevision: "project/\(ProjectName)/component/\(Name)"
		}
	}

	// Write a kustomization.yaml to the deploy directory for Kargo to edit.
	Artifacts: kargoKustomization: {
		artifact: "\(_OutPath)/kustomization.yaml"
		generators: [{
			kind:   "Resources"
			output: artifact
			resources: Kustomization: kargo: ks.#Kustomization & {
				resources: ["\(Name).gen.yaml"]
			}
		}]
	}
}

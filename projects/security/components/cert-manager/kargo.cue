@if(!NoKargo)
package holos

Parameters: {
	KargoProjectName: string @tag(KargoProjectName)
	KargoStageName:   string @tag(KargoStageName)
}

// Configure the ArgoCD Application to allow updates from Kargo.  Configure a
// stub kustomization.yaml artifact in the output directory for Kargo to edit.
Component: {
	_ArgoApplication: {
		metadata: annotations: "kargo.akuity.io/authorized-stage": "\(Parameters.KargoProjectName):\(Parameters.KargoStageName)"
	}
}

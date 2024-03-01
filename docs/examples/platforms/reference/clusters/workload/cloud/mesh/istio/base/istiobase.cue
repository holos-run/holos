package holos

#InputKeys: component: "istio-base"
#TargetNamespace: "istio-system"

#HelmChart & {
	namespace: #TargetNamespace
	chart: {
		name:    "base"
		version: "1.20.3"
		repository: {
			name: "istio"
			url:  "https://istio-release.storage.googleapis.com/charts"
		}
	}
	values: {
		global: {
			// Used to locate istiod.
			istioNamespace: #TargetNamespace
			// Switch the hub away from the default docker.io to avoid rate limits
			hub: "gcr.io/istio-release"
			// ImagePullSecrets for control plane ServiceAccount, list of secrets in the same namespace
      // to use for pulling any images in pods that reference this ServiceAccount.
      // Must be set for any cluster configured with private docker registry.
      imagePullSecrets: []
      istiod: enableAnalysis: false
      configValidation: true
      externalIstiod: false
      remotePilotAddress: ""
		}
		base: {
			// Include the CRDs in the helm template output
			enableCRDTemplates: true
			// Validation webhook configuration url
			// For example: https://$remotePilotAddress:15017/validate
			validationURL: ""
			// For istioctl usage to disable istio config crds in base
			enableIstioConfigCRDs: true
		}
		defaultRevision: "default"
	}
}

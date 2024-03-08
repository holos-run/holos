package holos

#IstioValues: {
	global: {
		// Used to locate istiod.
		istioNamespace: "istio-system"
		// Switch the hub away from the default docker.io to avoid rate limits
		hub: "gcr.io/istio-release"
		// ImagePullSecrets for control plane ServiceAccount, list of secrets in the same namespace
		// to use for pulling any images in pods that reference this ServiceAccount.
		// Must be set for any cluster configured with private docker registry.
		imagePullSecrets: []
		istiod: enableAnalysis: false
		configValidation:   true
		externalIstiod:     false
		remotePilotAddress: ""
	}
	base: {
		// holos includes crd templates with the --include-crds helm flag.
		enableCRDTemplates: false
		// Validation webhook configuration url
		// For example: https://$remotePilotAddress:15017/validate
		validationURL: ""
		// For istioctl usage to disable istio config crds in base
		enableIstioConfigCRDs: true
	}
	defaultRevision: "default"
}

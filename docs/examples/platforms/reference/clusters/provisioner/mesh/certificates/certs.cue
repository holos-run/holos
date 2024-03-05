package holos

// Provision all platform certificates.
#InputKeys: component: "certificates"

// Certificates usually go into the istio-system namespace, but they may go anywhere.
#TargetNamespace: "default"

// Depends on issuers
#DependsOn: _LetsEncrypt

#KubernetesObjects & {
	apiObjects: {
		for k, obj in #PlatformCerts {
			"\(obj.kind)": {
				"\(obj.metadata.namespace)/\(obj.metadata.name)": obj
			}
		}
	}
}

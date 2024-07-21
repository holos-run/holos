package holos

import (
	"encoding/json"
	core "github.com/holos-run/holos/api/core/v1alpha2"
	dto "github.com/holos-run/holos/service/gen/holos/object/v1alpha1:object"
	corev1 "k8s.io/api/core/v1"
	certv1 "cert-manager.io/certificate/v1"
	es "external-secrets.io/externalsecret/v1beta1"
)

// _PlatformConfig represents all of the data passed from holos to cue, used to
// carry the platform and project models.
_PlatformConfig:     dto.#PlatformConfig & json.Unmarshal(_PlatformConfigJSON)
_PlatformConfigJSON: string | *"{}" @tag(platform_config, type=string)

// #Cluster represents a single cluster in the platform.
#Cluster: {
	// name represents the name of the cluster.
	name: string
	// primary is true if the cluster is the sole primary cluster within the scope
	// of a fleet.
	primary: true | *false
}

_Clusters: #Clusters
// #Clusters defines the shape of _Clusters and clusters fields of other
// collections like #Fleet.clusters.
#Clusters: [Name=string]: #Cluster & {name: Name}

// _Fleets represents all the fleets in the platform.
_Fleets: #Fleets
// #Fleets defines the shape of _Fleets
#Fleets: [Name=string]: #Fleet & {name: Name}

// #Fleet represents a grouping of similar clusters.  A platform is usually
// composed of a workload fleet and a management fleet.
#Fleet: {
	name:     string
	clusters: #Clusters
}

// _Platform represents and provides a platform to holos for rendering.
_Platform: #Platform & {
	Name:  string @tag(platform_name, type=string)
	Model: _PlatformConfig.platformModel
}
// #Platform defines the shape of _Platform.
#Platform: {
	Name: string | *"holos"

	// Components represent the platform components to render.
	Components: [string]: core.#PlatformSpecComponent

	// Model represents the platform model from the web app form.
	Model: dto.#PlatformConfig.platformModel

	Output: core.#Platform & {
		metadata: name: Name

		spec: {
			// model represents the web form values provided by the user.
			model: Model
			components: [for c in Components {c}]
		}
	}
}

// _Namespaces represents all managed namespaces in the platform.
_Namespaces: #Namespaces
// #Namespaces defines the shape of _Namespaces.
#Namespaces: {
	[Name=string]: corev1.#Namespace & {
		metadata: name: Name
	}
}

// _Certificates represents all managed public facing tls certificates in the
// platform.
_Certificates: #Certificates
// #Certificates defines the shape of _Certificates
#Certificates: {
	[Name=string]: certv1.#Certificate & {
		metadata: name: Name
	}
}

// _Projects represents holos projects in the platform.
_Projects: #Projects
// #Projects defines the shape of _Projects
#Projects: [Name=string]: #Project & {
	metadata: name: Name
}

// #Project defines the shape of one project.
#Project: {
	metadata: name: string

	spec: {
		// namespaces represents the namespaces associated with this project.
		namespaces: #Namespaces
		// certificates represents the public tls certs associated with this project.
		certificates: #Certificates
	}
}

// #IngressCertificate defines a certificate for use by the ingress gateway.
#IngressCertificate: certv1.#Certificate & {
	metadata: name:      string
	metadata: namespace: string | *#IstioGatewaysNamespace
	spec: {
		commonName: string | *metadata.name
		secretName: metadata.name
		dnsNames: [...string] | *[commonName]
		issuerRef: kind: "ClusterIssuer"
		issuerRef: name: string | *"letsencrypt"
	}
}

// #ExternalSecret represents a typical external secret resource in the holos
// platform.  The default SecretStore in the same namespace is used.
#ExternalSecret: es.#ExternalSecret & {
	metadata: name: string
	spec: {
		target: name: metadata.name
		dataFrom: [{extract: {key: metadata.name}}]
		refreshInterval: "1h"
		secretStoreRef: kind: "SecretStore"
		secretStoreRef: name: "default"
	}
}

// #ExternalCert represents a tls Certificate managed in the management cluster and
// synced to the workload cluster using an ExternalSecret.
#ExternalCert: es.#ExternalSecret & {
	metadata: name:      string
	metadata: namespace: string | *#IstioGatewaysNamespace
	spec: {
		target: name: metadata.name
		target: template: type: "kubernetes.io/tls"
		target: creationPolicy: "Owner"
		target: deletionPolicy: "Retain"
		dataFrom: [
			{
				extract: {
					key:                metadata.name
					conversionStrategy: "Default"
					decodingStrategy:   "None"
					metadataPolicy:     "None"
				}
			},
		]
		refreshInterval: string | *"1h"
		secretStoreRef: kind: "SecretStore"
		secretStoreRef: name: string | *"default"
	}
}

// #IstioGatewaysNamespace represents the namespace where kubernetes Gateway API
// resources are deployed for istio.  This namespace was previously named
// "istio-ingress" when the istio Gateway API was used.
#IstioGatewaysNamespace: "istio-gateways"

// #Selector represents label selectors.
#Selector: [string]: matchLabels: {[string]: string}
_Selector: #Selector

// #AppInfo represents the data structure for an application deployed onto the
// platform.
#AppInfo: {
	metadata: {
		name:      string
		namespace: string
		labels: {[string]: string}
		annotations: {[string]: string}
	}

	spec: env:       string
	spec: component: string

	spec: region: hostname: string
	spec: global: hostname: string

	spec: dns: segments: {
		env: [] | [string]
		name: [] | [string]
		cluster: [] | [string]
		domain: [] | [string]
	}

	// The primary port for HTTPRoute
	spec: port: number

	spec: selector: matchLabels: {[string]: string}

	status: component: string
}

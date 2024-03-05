package holos

import "encoding/yaml"

let Name = "zitadel"
#InputKeys: component: Name

// Upstream helm chart doesn't specify the namespace field for all resources.
#Kustomization: spec: targetNamespace: #TargetNamespace

#HelmChart & {
	namespace: #TargetNamespace
	chart: {
		name:    Name
		version: "7.9.0"
		repository: {
			name: Name
			url:  "https://charts.zitadel.com"
		}
	}
	values: #Values

	apiObjects: {
		ExternalSecret: masterkey: #ExternalSecret & {
			_name: "zitadel-masterkey"
		}
		Certificate: zitadel: #Certificate & {
			metadata: name:      "crdb-zitadel-client"
			metadata: namespace: #TargetNamespace
			spec: {
				commonName: "zitadel"
				issuerRef: {
					group: "cert-manager.io"
					kind:  "Issuer"
					name:  "crdb-ca-issuer"
				}
				privateKey: algorithm: "RSA"
				privateKey: size:      2048
				renewBefore: "48h0m0s"
				secretName:  "cockroachdb-zitadel"
				subject: organizations: ["Cockroach"]
				usages: ["digital signature", "key encipherment", "client auth"]
			}
		}
		VirtualService: zitadel: #VirtualService & {
			metadata: name:      Name
			metadata: namespace: #TargetNamespace
			spec: hosts: ["login.\(#Platform.org.domain)"]
			spec: gateways: ["istio-ingress/default"]
			spec: http: [{route: [{destination: host: Name}]}]
		}
	}
}

// TODO: Generalize this common pattern of injecting the istio sidecar into a Deployment
let Patch = [{op: "add", path: "/spec/template/metadata/labels/sidecar.istio.io~1inject", value: "true"}]

#Kustomize: {
	patches: [
		{
			target: {
				group:   "apps"
				version: "v1"
				kind:    "Deployment"
				name:    Name
			}
			patch: yaml.Marshal(Patch)
		},
	]
}

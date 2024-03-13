package holos

import "encoding/yaml"

import "list"

let Name = "vault"
#InputKeys: component: Name
#InputKeys: project:   "core"
#TargetNamespace: "\(#InstancePrefix)-\(Name)"

let Vault = #OptionalServices[Name]

if Vault.enabled && list.Contains(Vault.clusterNames, #ClusterName) {
	#HelmChart & {
		namespace: #TargetNamespace
		chart: {
			name:    Name
			version: "0.25.0"
			repository: {
				name: "hashicorp"
				url:  "https://helm.releases.hashicorp.com"
			}
		}
		values: #Values

		apiObjects: {
			ExternalSecret: "gcpkms-creds":      _
			ExternalSecret: "vault-server-cert": _
			VirtualService: "\(Name)": {
				metadata: name:      Name
				metadata: namespace: #TargetNamespace
				spec: hosts: [for cert in Vault.certs {cert.spec.commonName}]
				spec: gateways: ["istio-ingress/\(Name)"]
				spec: http: [
					{
						route: [
							{
								destination: host: "\(Name)-active"
								destination: port: number: 8200
							},
						]
					},
				]
			}
		}
	}

	#Kustomize: {
		patches: [
			{
				target: {
					group:   "apps"
					version: "v1"
					kind:    "StatefulSet"
					name:    Name
				}
				patch: yaml.Marshal(EnvPatch)
			},
		]
	}

	let EnvPatch = [
		{
			op:    "test"
			path:  "/spec/template/spec/containers/0/env/4/name"
			value: "VAULT_ADDR"
		},
		{
			op:    "replace"
			path:  "/spec/template/spec/containers/0/env/4/value"
			value: "http://$(VAULT_K8S_POD_NAME):8200"
		},
	]
}

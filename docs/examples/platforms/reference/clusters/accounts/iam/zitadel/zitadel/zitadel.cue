package holos

import "encoding/yaml"

let Name = "zitadel"
#InputKeys: component: Name
#DependsOn: postgres:  _

// Upstream helm chart doesn't specify the namespace field for all resources.
#Kustomization: spec: targetNamespace: #TargetNamespace

#HelmChart & {
	namespace:   #TargetNamespace
	enableHooks: true
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
		ExternalSecret: "zitadel-masterkey": _
		VirtualService: "\(Name)": {
			metadata: name:      Name
			metadata: namespace: #TargetNamespace
			spec: hosts: ["login.\(#Platform.org.domain)"]
			spec: gateways: ["istio-ingress/default"]
			spec: http: [{route: [{destination: host: Name}]}]
		}
	}
}

// TODO: Generalize this common pattern of injecting the istio sidecar into a Deployment
let IstioInject = [{op: "add", path: "/spec/template/metadata/labels/sidecar.istio.io~1inject", value: "true"}]

_PGBouncer: "pgbouncer"

let DatabaseCACertPatch = [
	{
		op:   "add"
		path: "/spec/template/spec/volumes/-"
		value: {
			name: _PGBouncer
			secret: {
				secretName: "\(_DBName)-pgbouncer"
				items: [{key: "pgbouncer-frontend.ca-roots", path: "ca.crt"}]
			}
		}
	},
	{
		op:   "add"
		path: "/spec/template/spec/containers/0/volumeMounts/-"
		value: {
			name:      _PGBouncer
			mountPath: "/" + _PGBouncer
		}
	},
]

#Kustomize: {
	patches: [
		{
			target: {
				group:   "apps"
				version: "v1"
				kind:    "Deployment"
				name:    Name
			}
			patch: yaml.Marshal(IstioInject)
		},
		{
			target: {
				group:   "apps"
				version: "v1"
				kind:    "Deployment"
				name:    Name
			}
			patch: yaml.Marshal(DatabaseCACertPatch)
		},
		{
			target: {
				group:   "batch"
				version: "v1"
				kind:    "Job"
				name:    "\(Name)-init"
			}
			patch: yaml.Marshal(DatabaseCACertPatch)
		},
		{
			target: {
				group:   "batch"
				version: "v1"
				kind:    "Job"
				name:    "\(Name)-setup"
			}
			patch: yaml.Marshal(DatabaseCACertPatch)
		},
	]
}

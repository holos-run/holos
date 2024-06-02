package holos

import (
	"encoding/yaml"
	kc "sigs.k8s.io/kustomize/api/types"
)

// Produce a helm chart build plan.
(#Helm & Chart).Output

let Cluster = _Clusters[_ClusterName]

let Chart = {
	Name:      "zitadel-server"
	Version:   "7.11.0"
	Namespace: _ZitadelNamespace

	Chart: chart: name: "zitadel"
	Chart: enableHooks: true

	Repo: name: "zitadel"
	Repo: url:  "https://charts.zitadel.com"

	Values: #Values

	Resources: [_]: [_]: metadata: namespace: Namespace
	Resources: ExternalSecret: "zitadel-masterkey": #ExternalSecret
	// Grant the Gateway namespace the ability to refer to the backend httpbin
	// service in HTTPRoutes.
	Resources: ReferenceGrant: (#IstioGatewaysNamespace): #ReferenceGrant

	EnableKustomizePostProcessor: true
	// Force all resources into the zitadel namespace, some resources in the helm
	// chart do not specify the namespace so they may get mis-applied depending on
	// the kubectl (client-go) context.
	KustomizeFiles: "kustomization.yaml": namespace: Namespace

	KustomizePatches: {
		mesh: {
			target: name:    "zitadel"
			target: kind:    "Deployment"
			target: version: "v1"
			target: group:   "apps"
			patch: yaml.Marshal(IstioInject)
		}

		deploymentCA: CAPatch & {
			target: name:  "zitadel"
			target: kind:  "Deployment"
			target: group: "apps"
		}

		initJob: CAPatch & {
			target: name:  "zitadel-init"
			target: kind:  "Job"
			target: group: "batch"
		}

		setupJob: CAPatch & {
			target: name:  "zitadel-setup"
			target: kind:  "Job"
			target: group: "batch"
		}

		testDisable: {
			target: name:    "zitadel-test-connection"
			target: kind:    "Pod"
			target: version: "v1"
			patch: yaml.Marshal(DisableFluxPatch)
		}

		if Cluster.primary == false {
			fluxDisable: {
				target: {
					group:   "apps"
					version: "v1"
					kind:    "Deployment"
					name:    Name
				}
				patch: yaml.Marshal(DisableFluxPatch)
			}

			initDisable: {
				target: {
					group:   "batch"
					version: "v1"
					kind:    "Job"
					name:    "\(Name)-init"
				}
				patch: yaml.Marshal(DisableFluxPatch)
			}

			setupDisable: {
				target: {
					group:   "batch"
					version: "v1"
					kind:    "Job"
					name:    "\(Name)-setup"
				}
				patch: yaml.Marshal(DisableFluxPatch)
			}
		}
	}
}

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

let CAPatch = kc.#Patch & {
	target: {
		group:   "apps" | "batch"
		version: "v1"
		kind:    "Job" | "Deployment"
		name:    string
	}
	patch: yaml.Marshal(DatabaseCACertPatch)
}

let DisableFluxPatch = [{op: "replace", path: "/metadata/annotations/kustomize.toolkit.fluxcd.io~1reconcile", value: "disabled"}]

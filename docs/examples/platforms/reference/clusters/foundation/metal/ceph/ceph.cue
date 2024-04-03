package holos

// Manage Ceph CSI to provide PersistentVolumeClaims to a cluster.

#TargetNamespace: "ceph-system"

#SecretName: "\(#ClusterName)-ceph-csi-rbd"

#Kustomization: spec: targetNamespace: "ceph-system"

spec: components: HelmChartList: [
	#HelmChart & {
		_dependsOn: "prod-secrets-stores": _

		metadata: name: "prod-metal-ceph"

		namespace: #TargetNamespace
		chart: {
			name:    "ceph-csi-rbd"
			version: "3.10.2"
			repository: {
				name: "ceph-csi"
				url:  "https://ceph.github.io/csi-charts"
			}
		}
		_values:      #ChartValues
		apiObjectMap: OBJECTS.apiObjectMap
	},
]

let OBJECTS = #APIObjects & {
	apiObjects: {
		ExternalSecret: "\(#SecretName)": _
	}
}

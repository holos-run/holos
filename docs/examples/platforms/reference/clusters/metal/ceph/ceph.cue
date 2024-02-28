package holos

// Manage Ceph CSI to provide PersistentVolumeClaims to a cluster.

#TargetNamespace: "ceph-system"

#InputKeys: {
	project:   "metal"
	service:   "ceph"
	component: "ceph"
}

#Kustomization: spec: {
	dependsOn: [{name: "prod-secrets-namespaces"}]
	targetNamespace: #TargetNamespace
}

#HelmChart & {
	namespace: #TargetNamespace
	chart: {
		name:    "ceph-csi-rbd"
		version: "3.10.2"
		repository: {
			name: "ceph-csi"
			url:  "https://ceph.github.io/csi-charts"
		}
	}
	objects: [#ExternalSecret & { _name: "ceph-csi-rbd" }]
}

package holos

#Input: {
	config : {
		  // (required) String representing a Ceph cluster to provision storage from.
		  // Should be unique across all Ceph clusters in use for provisioning,
		  // cannot be greater than 36 bytes in length, and should remain immutable for
		  // the lifetime of the StorageClass in use.
			clusterID: string
			// (required) []String list of ceph monitor "address:port" values.
			monitors: [...string]
	}
}

// Imported from https://github.com/holos-run/holos-infra/blob/0ae58858f5583d25fa7543e47b5f5e9f0b2f3c83/components/core/metal/ceph-csi-rbd/values.holos.yaml

#ChartValues: {
	// Necessary for Talos see https://github.com/siderolabs/talos/discussions/8163
	selinuxMount: false

	csiConfig: [#Input.config]

	storageClass: {
		annotations: "storageclass.kubernetes.io/is-default-class": "true"

		// Specifies whether the storageclass should be created
		create: true
		name:   "ceph-ssd"

		// (optional) Prefix to use for naming RBD images.
		// If omitted, defaults to "csi-vol-".
		// NOTE: Set this to a cluster specific value, e.g. vol-k1-
		volumeNamePrefix: "vol-\(#InputKeys.cluster)-"

		// (required) String representing a Ceph cluster to provision storage from.
		// Should be unique across all Ceph clusters in use for provisioning,
		// cannot be greater than 36 bytes in length, and should remain immutable for
		// the lifetime of the StorageClass in use.
		clusterID: #Input.config.clusterID

		// (optional) If you want to use erasure coded pool with RBD, you need to
		// create two pools. one erasure coded and one replicated.
		// You need to specify the replicated pool here in the `pool` parameter, it is
		// used for the metadata of the images.
		// The erasure coded pool must be set as the `dataPool` parameter below.
		// dataPool: <ec-data-pool>
		dataPool: ""

		// (required) Ceph pool into which the RBD image shall be created
		// eg: pool: replicapool
		pool: "k8s-dev"

		// (optional) RBD image features, CSI creates image with image-format 2 CSI
		// RBD currently supports `layering`, `journaling`, `exclusive-lock`,
		// `object-map`, `fast-diff`, `deep-flatten` features.
		// Refer https://docs.ceph.com/en/latest/rbd/rbd-config-ref/#image-features
		// for image feature dependencies.
		// imageFeatures: layering,journaling,exclusive-lock,object-map,fast-diff
		imageFeatures: "layering"

		// (optional) Specifies whether to try other mounters in case if the current
		// mounter fails to mount the rbd image for any reason. True means fallback
		// to next mounter, default is set to false.
		// Note: tryOtherMounters is currently useful to fallback from krbd to rbd-nbd
		// in case if any of the specified imageFeatures is not supported by krbd
		// driver on node scheduled for application pod launch, but in the future this
		// should work with any mounter type.
		// tryOtherMounters: false
		// (optional) uncomment the following to use rbd-nbd as mounter
		// on supported nodes
		// mounter: rbd-nbd
		mounter: ""

		// (optional) ceph client log location, eg: rbd-nbd
		// By default host-path /var/log/ceph of node is bind-mounted into
		// csi-rbdplugin pod at /var/log/ceph mount path. This is to configure
		// target bindmount path used inside container for ceph clients logging.
		// See docs/rbd-nbd.md for available configuration options.
		// cephLogDir: /var/log/ceph
		cephLogDir: ""

		// (optional) ceph client log strategy
		// By default, log file belonging to a particular volume will be deleted
		// on unmap, but you can choose to just compress instead of deleting it
		// or even preserve the log file in text format as it is.
		// Available options `remove` or `compress` or `preserve`
		// cephLogStrategy: remove
		cephLogStrategy: ""

		// (optional) Instruct the plugin it has to encrypt the volume
		// By default it is disabled. Valid values are "true" or "false".
		// A string is expected here, i.e. "true", not true.
		// encrypted: "true"
		encrypted: ""

		// (optional) Use external key management system for encryption passphrases by
		// specifying a unique ID matching KMS ConfigMap. The ID is only used for
		// correlation to configmap entry.
		encryptionKMSID: ""

		// Add topology constrained pools configuration, if topology based pools
		// are setup, and topology constrained provisioning is required.
		// For further information read TODO<doc>
		// topologyConstrainedPools: |
		//   [{"poolName":"pool0",
		//     "dataPool":"ec-pool0" # optional, erasure-coded pool for data
		//     "domainSegments":[
		//       {"domainLabel":"region","value":"east"},
		//       {"domainLabel":"zone","value":"zone1"}]},
		//    {"poolName":"pool1",
		//     "dataPool":"ec-pool1" # optional, erasure-coded pool for data
		//     "domainSegments":[
		//       {"domainLabel":"region","value":"east"},
		//       {"domainLabel":"zone","value":"zone2"}]},
		//    {"poolName":"pool2",
		//     "dataPool":"ec-pool2" # optional, erasure-coded pool for data
		//     "domainSegments":[
		//       {"domainLabel":"region","value":"west"},
		//       {"domainLabel":"zone","value":"zone1"}]}
		//   ]
		topologyConstrainedPools: []

		// (optional) mapOptions is a comma-separated list of map options.
		// For krbd options refer
		// https://docs.ceph.com/docs/master/man/8/rbd/#kernel-rbd-krbd-options
		// For nbd options refer
		// https://docs.ceph.com/docs/master/man/8/rbd-nbd/#options
		// Format:
		// mapOptions: "<mounter>:op1,op2;<mounter>:op1,op2"
		// An empty mounter field is treated as krbd type for compatibility.
		// eg:
		// mapOptions: "krbd:lock_on_read,queue_depth=1024;nbd:try-netlink"
		mapOptions: ""

		// (optional) unmapOptions is a comma-separated list of unmap options.
		// For krbd options refer
		// https://docs.ceph.com/docs/master/man/8/rbd/#kernel-rbd-krbd-options
		// For nbd options refer
		// https://docs.ceph.com/docs/master/man/8/rbd-nbd/#options
		// Format:
		// unmapOptions: "<mounter>:op1,op2;<mounter>:op1,op2"
		// An empty mounter field is treated as krbd type for compatibility.
		// eg:
		// unmapOptions: "krbd:force;nbd:force"
		unmapOptions: ""

		// The secrets have to contain Ceph credentials with required access
		// to the 'pool'.
		provisionerSecret: "csi-rbd-secret"
		// If Namespaces are left empty, the secrets are assumed to be in the
		// Release namespace.
		provisionerSecretNamespace:      ""
		controllerExpandSecret:          "csi-rbd-secret"
		controllerExpandSecretNamespace: ""
		nodeStageSecret:                 "csi-rbd-secret"
		nodeStageSecretNamespace:        ""
		// Specify the filesystem type of the volume. If not specified,
		// csi-provisioner will set default as `ext4`.
		fstype:               "ext4"
		reclaimPolicy:        "Delete"
		allowVolumeExpansion: true
		mountOptions: []
	}

	secret: {
		// Specifies whether the secret should be created
		create: false
		name:   "csi-rbd-secret"
		// Key values correspond to a user name and its key, as defined in the
		// ceph cluster. User ID should have required access to the 'pool'
		// specified in the storage class
		userID:  "admin"
		userKey: "$(ceph auth get-key client.admin)"
		// Encryption passphrase
		encryptionPassphrase: "$(python -c 'import secrets; print(secrets.token_hex(32));')"
	}
}

package holos

// If you are using k3d with the default Flannel CNI, you must append some
// values to your installation command, as k3d uses nonstandard locations for
// CNI configuration and binaries.
//
// See https://istio.io/latest/docs/ambient/install/platform-prerequisites/#k3d
#Istio: Values: cni: {
	cniConfDir: "/var/lib/rancher/k3s/agent/etc/cni/net.d"
	cniBinDir:  "/bin"
}

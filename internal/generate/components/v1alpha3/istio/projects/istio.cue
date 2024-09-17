package holos

// #Istio represents platform wide configuration
#Istio: {
	Version: "1.23.1"
	System: Namespace: "istio-system"

	// Constrain Helm values for safer, easier upgrades and consistency across
	// platform components.
	Values: global: istioNamespace: System.Namespace

	// Configure ambient mode
	Values: profile: "ambient"
}

// Register the Namespaces
#Namespaces: (#Istio.System.Namespace): _

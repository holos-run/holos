package holos

// #Platform represents the user supplied platform configuration.
#Platform: {
	#HolosTypeMeta
	kind:       "Platform"
	apiVersion: "app.holos.run/v1alpha1"
	metadata:   #HolosObjectMeta
	spec:       #PlatformSpec
	holos:      #Holos
}

// #Holos represents the holos reserved field in the #Platform schema defined by the holos development team.
#Holos: {
	// flags represents config values provided by holos command line flags.
	flags: {
		// cluster represents the holos render --cluster-name flag.
		cluster: string @tag(cluster, type=string)
	}
}

// #PlatformSpec represents configuration values defined by the platform
// designer.  Config values are organized by section, then simple strings for
// each section.
#PlatformSpec: {
	config: [string]: _
	config: user:     #UserDefinedConfig
}

// #PlatformUserConfig represents configuration fields and values defined by the
// user.
#UserDefinedConfig: {
	sections: [string]: fields: [string]: _
}

// #PlatformConfig represents the platform config data returned from the Holos API.  Useful for cue vet.
#PlatformConfig: {
	platform: spec: #PlatformSpec
}

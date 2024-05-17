package v1alpha1

// #Platform represents a platform to manage.  Holos manages a platform by
// rendering platform components and applying the configuration to clusters as
// defined by the platform resource.
#Platform: {
	#TypeMeta
	apiVersion: #APIVersion
	kind:       "Platform"
	metadata:   #ObjectMeta
	spec: {
		// model represents the user defined platform model, which is produced and
		// defined by the user supplied form.
		model: {...}
	}
}

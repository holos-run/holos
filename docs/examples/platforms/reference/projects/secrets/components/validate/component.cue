package holos

objects: [
	#SecretStore,
	#ExternalSecret & {_name: "validate"},
]

{} & #KubernetesObjects

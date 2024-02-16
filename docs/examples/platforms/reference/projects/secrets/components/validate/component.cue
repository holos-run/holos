package holos

#Kustomization: spec: dependsOn: [{name: #InstancePrefix + "-eso"}]

objects: [
	#SecretStore,
	#ExternalSecret & {
		_name: "validate"
		spec: dataFrom: [{extract: key: "ns/" + #TargetNamespace + "/test"}]
	},
]

{} & #KubernetesObjects

package holos

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"encoding/yaml"
)

// #Name defines the name: string key value pair used all over the place.
#Name: name: string

// #NamespaceMeta defines standard metadata for namespaces.
// Refer to https://kubernetes.io/docs/reference/labels-annotations-taints/#kubernetes-io-metadata-name
#NamespaceMeta: {
	metadata: {
		name: string
		labels: "kubernetes.io/metadata.name": name
	}
	...
}

// Kubernetes API Objects
#Namespace: corev1.#Namespace & #NamespaceMeta
#ConfigMap: corev1.#ConfigMap

// #InputKeys defines the set of cue tags required to build a cue holos component. The values are used as lookup keys into the _Platform data.
#InputKeys: {
	// cluster is usually the only key necessary when working with a component on the command line.
	cluster: string @tag(cluster, type=string)
	// stage is usually set by the platform or project.
	stage: string @tag(stage, type=string)
	// project is usually set by the platform or project.
	project: string @tag(project, type=string)
	// service is usually set by the component.
	service: string @tag(service, type=string)
}

// #Platform defines the primary lookup table for the platform.  Lookup keys should be limited to those defined in #KeyTags.
#Platform: {
	// org holds user defined values scoped organization wide.  A platform has one and only one organization.
	org: {
		name: string
		domain: string
	}
	clusters: [ID=_]: {
		name: string & ID
		region?: string
	}
	stages: [ID=_]: {
		name: string & ID
		environments: [...#Name]
	}
	projects: [ID=_]: {
		name: string & ID
	}
	services: [ID=_]: {
		name: string & ID
	}
}
// _PlatformData stores the values of the primary lookup table.
_Platform: #Platform

// #ComponentOutput is the output schema of a single component.
#ComponentOutput: {
	// objects holds a list of the kubernetes api objects to configure.
	objects: [...metav1.#TypeMeta] | *[]
	// out holds the rendered yaml stream of kubernetes api objects.
	out: yaml.MarshalStream(objects)
	// platform returns the platform data structure for visibility / troubleshooting.
	platform: _Platform
	// debug returns arbitrary debug output.
	debug?: _
}

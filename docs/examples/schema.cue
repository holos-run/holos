package holos

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ksv1 "kustomize.toolkit.fluxcd.io/kustomization/v1"
	corev1 "k8s.io/api/core/v1"
	"encoding/yaml"
)

_apiVersion: "holos.run/v1alpha1"

// #Name defines the name: string key value pair used all over the place.
#Name: name: string

// #InstanceName is the name of the holos component instance being managed varying by stage, project, and component names.
#InstanceName: "\(#InputKeys.stage)-\(#InputKeys.project)-\(#InputKeys.component)"

// #NamespaceMeta defines standard metadata for namespaces.
// Refer to https://kubernetes.io/docs/reference/labels-annotations-taints/#kubernetes-io-metadata-name
#NamespaceMeta: {
	metadata: {
		name: string
		labels: "kubernetes.io/metadata.name": name
	}
	...
}

// #TargetNamespace is the target namespace for a holos component.
#TargetNamespace: string

// Kubernetes API Objects
#Namespace: corev1.#Namespace & #NamespaceMeta
#ConfigMap: corev1.#ConfigMap
#Kustomization: ksv1.#Kustomization & {
	metadata: {
		name: #InstanceName,
		namespace: string | *"flux-system",
	}
	spec: ksv1.#KustomizationSpec & {
		interval: string | *"30m0s"
		path: string | *"deploy/clusters/\(#InputKeys.cluster)/components/\(#InstanceName)"
		prune: bool | *true
		retryInterval: string | *"2m0s"
		sourceRef: {
			kind: string | *"GitRepository"
			name: string | *"flux-system"
		}
		timeout: string | *"3m0s"
		wait: bool | *true
	}
}


// #InputKeys defines the set of cue tags required to build a cue holos component. The values are used as lookup keys into the _Platform data.
#InputKeys: {
	// cluster is usually the only key necessary when working with a component on the command line.
	cluster: string @tag(cluster, type=string)
	// stage is usually set by the platform or project.
	stage: *"prod" | string @tag(stage, type=string)
	// project is usually set by the platform or project.
	project: string @tag(project, type=string)
	// service is usually set by the component.
	service: string @tag(service, type=string)
	// component is the name of the component
	component: string @tag(component, type=string)
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

// #OutputTypeMeta is shared among all output types
#OutputTypeMeta: {
	// apiVersion is the output api version
	apiVersion: _apiVersion
	// kind is a discriminator of the type of output
	kind: #PlatformSpec.kind | #KubernetesObjects.kind | #HelmChart.kind
	// name holds a unique name suitable for a filename
	metadata: name: string
	// contentType is the standard MIME type indicating the content type of the content field
	contentType: *"application/yaml" | "application/json"
	// content holds the content text output
	content: string | *""
	// debug returns arbitrary debug output.
	debug?: _
}

// #KubernetesObjectOutput is the output schema of a single component.
#KubernetesObjects: {
	#OutputTypeMeta
	// kind KubernetesObjects provides a yaml text stream of kubernetes api objects in the out field.
	kind: "KubernetesObjects"
	// objects holds a list of the kubernetes api objects to configure.
	objects: [...metav1.#TypeMeta] | *[]
	// out holds the rendered yaml text stream of kubernetes api objects.
	content: yaml.MarshalStream(objects)
	// ksObjects holds the flux Kustomization objects for gitops
	ksObjects: [...#Kustomization] | *[]
	// ksContent is the yaml representation of kustomization
	ksContent: yaml.MarshalStream(ksObjects)
	// platform returns the platform data structure for visibility / troubleshooting.
	platform: _Platform
}

// #Chart defines an upstream helm chart
#Chart: {
	name: string
	version: string
	repository: {
		name: string
		url: string
	}
}

// #HelmChart is a holos component which produces kubernetes api objects from cue values provided to the helm template command.
#HelmChart: {
	#OutputTypeMeta
	kind: "HelmChart"
	// ksObjects holds the flux Kustomization objects for gitops.
	ksObjects: [...#Kustomization] | *[#Kustomization]
	// ksContent is the yaml representation of kustomization.
	ksContent: yaml.MarshalStream(ksObjects)
	// namespace defines the value passed to the helm --namespace flag
	namespace: #TargetNamespace
	// chart defines the upstream helm chart to process.
	chart: #Chart
	// values represents the helm values to provide to the chart.
	values: {...}
	// valuesContent holds the values yaml
	valuesContent: yaml.Marshal(values)
	// platform returns the platform data structure for visibility / troubleshooting.
	platform: _Platform
	// instance returns the key values of the holos component instance.
	instance: #InputKeys
}

// #PlatformSpec is the output schema of a platform specification.
#PlatformSpec: {
	#OutputTypeMeta
	kind: "PlatformSpec"
}

#Output: #PlatformSpec | #KubernetesObjects | #HelmChart

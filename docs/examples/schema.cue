package holos

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ksv1 "kustomize.toolkit.fluxcd.io/kustomization/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	batchv1 "k8s.io/api/batch/v1"
	es "external-secrets.io/externalsecret/v1beta1"
	ss "external-secrets.io/secretstore/v1beta1"
	cm "cert-manager.io/clusterissuer/v1"
	"encoding/yaml"
)

// _apiVersion is the version of this schema.  Defines the interface between CUE output and the holos cli.
_apiVersion: "holos.run/v1alpha1"

// #ClusterName is the cluster name for cluster scoped resources.
#ClusterName: #InputKeys.cluster

// #StageName is prod, dev, stage, etc...  Usually prod for platform components.
#StageName: #InputKeys.stage

// #CollectionName is the preferred handle to the collection element of the instance name.  A collection name mapes to an "application name" as described in the kubernetes recommended labels documentation.  Refer to https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
#CollectionName: #InputKeys.project

// #ComponentName is the name of the holos component.
#ComponentName: #InputKeys.component

// #InstanceName is the name of the holos component instance being managed varying by stage, project, and component names.
#InstanceName: "\(#StageName)-\(#CollectionName)-\(#ComponentName)"

// #InstancePrefix is the stage and project without the component name.  Useful for dependency management among multiple components for a project stage.
#InstancePrefix: "\(#StageName)-\(#CollectionName)"

// #TargetNamespace is the target namespace for a holos component.
#TargetNamespace: string

// #CommonLabels are mixed into every kubernetes api object.
#CommonLabels: {
	"holos.run/stage.name":        #StageName
	"holos.run/project.name":      #CollectionName
	"holos.run/component.name":    #ComponentName
	"app.kubernetes.io/part-of":   #StageName
	"app.kubernetes.io/name":      #CollectionName
	"app.kubernetes.io/component": #ComponentName
	"app.kubernetes.io/instance":  #InstanceName
	...
}

#ClusterObject: {
	metadata: metav1.#ObjectMeta & {
		labels: #CommonLabels
	}
	...
}

#NamespaceObject: #ClusterObject & {
	metadata: namespace: string
}

// Kubernetes API Objects
#Namespace: corev1.#Namespace & #ClusterObject & {
	metadata: {
		name: string
		labels: "kubernetes.io/metadata.name": name
	}
}
#ClusterRole:        #ClusterObject & rbacv1.#ClusterRole
#ClusterRoleBinding: #ClusterObject & rbacv1.#ClusterRoleBinding
#ClusterIssuer: #ClusterObject & cm.#ClusterIssuer & {...}
#Role:           #NamespaceObject & rbacv1.#Role
#RoleBinding:    #NamespaceObject & rbacv1.#RoleBinding
#ConfigMap:      #NamespaceObject & corev1.#ConfigMap
#ServiceAccount: #NamespaceObject & corev1.#ServiceAccount
#Pod:            #NamespaceObject & corev1.#Pod
#Job:            #NamespaceObject & batchv1.#Job
#CronJob:        #NamespaceObject & batchv1.#CronJob

// Flux Kustomization CRDs
#Kustomization: #NamespaceObject & ksv1.#Kustomization & {
	metadata: {
		name:      #InstanceName
		namespace: string | *"flux-system"
	}
	spec: ksv1.#KustomizationSpec & {
		interval:      string | *"30m0s"
		path:          string | *"deploy/clusters/\(#InputKeys.cluster)/components/\(#InstanceName)"
		prune:         bool | *true
		retryInterval: string | *"2m0s"
		sourceRef: {
			kind: string | *"GitRepository"
			name: string | *"flux-system"
		}
		suspend?:         bool
		targetNamespace?: string
		timeout:          string | *"3m0s"
		wait:             bool | *true
	}
}

// External Secrets CRDs
#ExternalSecret: #NamespaceObject & es.#ExternalSecret & {
	_name: string
	metadata: {
		name:      _name
		namespace: #TargetNamespace
	}
	spec: {
		refreshInterval: string | *"1h"
		secretStoreRef: {
			kind: string | *"SecretStore"
			name: string | *"default"
		}
		target: {
			name:           _name
			creationPolicy: string | *"Owner"
			deletionPolicy: string | *"Retain"
		}
		// Copy fields 1:1 from external Secret to target Secret.
		dataFrom: [{extract: key: _name}]
	}
}

#SecretStore: #NamespaceObject & ss.#SecretStore & {
	metadata: {
		name:      string | *"default"
		namespace: #TargetNamespace
	}
	spec: provider: {
		kubernetes: {
			remoteNamespace: #TargetNamespace
			auth: token: bearerToken: {
				name: string | *"eso-reader"
				key:  string | *"token"
			}
			server: {
				caBundle: #InputKeys.provisionerCABundle
				url:      #InputKeys.provisionerURL
			}
		}
	}
}

// #InputKeys defines the set of cue tags required to build a cue holos component. The values are used as lookup keys into the #Platform data.
#InputKeys: {
	// cluster is usually the only key necessary when working with a component on the command line.
	cluster: string @tag(cluster, type=string)
	// stage is usually set by the platform or project.
	stage: *"prod" | string @tag(stage, type=string)
	// service is usually set by the component.
	service: *component | string @tag(service, type=string)
	// component is the name of the component
	component: string @tag(component, type=string)

	// GCP Project Info used for the Provisioner Cluster
	gcpProjectID:     string @tag(gcpProjectID, type=string)
	gcpProjectNumber: int    @tag(gcpProjectNumber, type=int)

	// Same as cluster certificate-authority-data field in ~/.holos/kubeconfig.provisioner
	provisionerCABundle: string @tag(provisionerCABundle, type=string)
	// Same as the cluster server field in ~/.holos/kubeconfig.provisioner
	provisionerURL: string @tag(provisionerURL, type=string)
}

// #Platform defines the primary lookup table for the platform.  Lookup keys should be limited to those defined in #KeyTags.
#Platform: {
	// org holds user defined values scoped organization wide.  A platform has one and only one organization.
	org: {
		name:   string
		domain: string
		contact: email:    string
		cloudflare: email: string
	}
	clusters: [ID=_]: {
		name:    string & ID
		region?: string
	}
	stages: [ID=_]: {
		name: string & ID
		environments: [...{name: string}]
	}
	projects: [ID=_]: {
		name: string & ID
	}
	services: [ID=_]: {
		name: string & ID
	}
}

// #APIObjects is the output type for api objects produced by cue.  A map is used to aid debugging and clarity.
#APIObjects: {
	// apiObjects holds each the api objects produced by cue.
	apiObjects: {
		[Kind=_]: {
			[Name=_]: metav1.#TypeMeta & {
				kind: Kind
			}
		}
	}

	// apiObjectsContent holds the marshalled representation of apiObjects
	apiObjectMap: {
		for kind, v in apiObjects {
			"\(kind)": {
				for name, obj in v {
					"\(name)": yaml.Marshal(obj)
				}
			}
		}
		...
	}
}

// #OutputTypeMeta is shared among all output types
#OutputTypeMeta: {
	// apiVersion is the output api version
	apiVersion: _apiVersion
	// kind is a discriminator of the type of output
	kind: #PlatformSpec.kind | #KubernetesObjects.kind | #HelmChart.kind | #NoOutput.kind
	// name holds a unique name suitable for a filename
	metadata: name: string
	// debug returns arbitrary debug output.
	debug?: _
}

#NoOutput: {
	#OutputTypeMeta
	kind: string | *"Skip"
	metadata: name: string | *"skipped"
}

// #KubernetesObjectOutput is the output schema of a single component.
#KubernetesObjects: {
	#OutputTypeMeta
	#APIObjects
	kind: "KubernetesObjects"
	metadata: name: #InstanceName
	// ksObjects holds the flux Kustomization objects for gitops
	ksObjects: [...#Kustomization] | *[#Kustomization]
	// ksContent is the yaml representation of kustomization
	ksContent: yaml.Marshal(#Kustomization)
	// platform returns the platform data structure for visibility / troubleshooting.
	platform: #Platform
}

// #Chart defines an upstream helm chart
#Chart: {
	name:    string
	version: string
	repository: {
		name: string
		url:  string
	}
}

// #ChartValues represent the values provided to a helm chart.  Existing values may be imorted using cue import values.yaml -p holos then wrapping the values.cue content in #Values: {}
#ChartValues: {...}

// #HelmChart is a holos component which produces kubernetes api objects from cue values provided to the helm template command.
#HelmChart: {
	#OutputTypeMeta
	#APIObjects
	kind: "HelmChart"
	metadata: name: #InstanceName
	// ksObjects holds the flux Kustomization objects for gitops.
	ksObjects: [...#Kustomization] | *[#Kustomization]
	// ksContent is the yaml representation of kustomization.
	ksContent: yaml.MarshalStream(ksObjects)
	// namespace defines the value passed to the helm --namespace flag
	namespace: #TargetNamespace
	// chart defines the upstream helm chart to process.
	chart: #Chart
	// values represents the helm values to provide to the chart.
	values: #ChartValues
	// valuesContent holds the values yaml
	valuesContent: yaml.Marshal(values)
	// platform returns the platform data structure for visibility / troubleshooting.
	platform: #Platform
	// instance returns the key values of the holos component instance.
	instance: #InputKeys
}

// #PlatformSpec is the output schema of a platform specification.
#PlatformSpec: {
	#OutputTypeMeta
	kind: "PlatformSpec"
}

// #SecretName is the name of a Secret, ususally coupling a Deployment to an ExternalSecret
#SecretName: string

// By default, render kind: Skipped so holos knows to skip over intermediate cue files.
// This enables the use of holos render ./foo/bar/baz/... when bar contains intermediary constraints which are not complete components.
// Holos skips over these intermediary cue instances.
{} & #NoOutput

package holos

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ksv1 "kustomize.toolkit.fluxcd.io/kustomization/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	batchv1 "k8s.io/api/batch/v1"
	es "external-secrets.io/externalsecret/v1beta1"
	ss "external-secrets.io/secretstore/v1beta1"
	is "cert-manager.io/issuer/v1"
	ci "cert-manager.io/clusterissuer/v1"
	crt "cert-manager.io/certificate/v1"
	gw "networking.istio.io/gateway/v1beta1"
	vs "networking.istio.io/virtualservice/v1beta1"
	kc "sigs.k8s.io/kustomize/api/types"
	pg "postgres-operator.crunchydata.com/postgrescluster/v1beta1"
	"encoding/yaml"
)

let ResourcesFile = "resources.yaml"

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

// #SelectorLabels are mixed into selectors.
#SelectorLabels: {
	"holos.run/stage.name":     #StageName
	"holos.run/project.name":   #CollectionName
	"holos.run/component.name": #ComponentName
	...
}

// #CommonLabels are mixed into every kubernetes api object.
#CommonLabels: {
	#SelectorLabels
	"app.kubernetes.io/part-of":   #StageName
	"app.kubernetes.io/name":      #CollectionName
	"app.kubernetes.io/component": #ComponentName
	"app.kubernetes.io/instance":  #InstanceName
	...
}

#ClusterObject: {
	_description: string | *""
	metadata: metav1.#ObjectMeta & {
		labels: #CommonLabels
		annotations: #Description & {
			_Description: _description
			...
		}
	}
	...
}

#Description: {
	_Description:            string | *""
	"holos.run/description": _Description
	...
}

#NamespaceObject: #ClusterObject & {
	metadata: name:      string
	metadata: namespace: string
	...
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
#ClusterIssuer: #ClusterObject & ci.#ClusterIssuer & {...}

#Issuer:          #NamespaceObject & is.#Issuer
#Role:            #NamespaceObject & rbacv1.#Role
#RoleBinding:     #NamespaceObject & rbacv1.#RoleBinding
#ConfigMap:       #NamespaceObject & corev1.#ConfigMap
#ServiceAccount:  #NamespaceObject & corev1.#ServiceAccount
#Pod:             #NamespaceObject & corev1.#Pod
#Service:         #NamespaceObject & corev1.#Service
#Job:             #NamespaceObject & batchv1.#Job
#CronJob:         #NamespaceObject & batchv1.#CronJob
#Deployment:      #NamespaceObject & appsv1.#Deployment
#Gateway:         #NamespaceObject & gw.#Gateway
#VirtualService:  #NamespaceObject & vs.#VirtualService
#Certificate:     #NamespaceObject & crt.#Certificate
#PostgresCluster: #NamespaceObject & pg.#PostgresCluster

// #HTTP01Cert defines a http01 certificate.
#HTTP01Cert: {
	_name:      string
	_secret:    string | *_name
	SecretName: _secret
	Host:       _name + "." + #ClusterDomain
	object: #Certificate & {
		metadata: {
			name:      _secret
			namespace: string | *#TargetNamespace
		}
		spec: {
			commonName: Host
			dnsNames: [Host]
			secretName: _secret
			issuerRef: kind: "ClusterIssuer"
			issuerRef: name: "letsencrypt"
		}
	}
}

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
		// wait performs health checks for all reconciled resources. If set to true, .spec.healthChecks is ignored.
		// Setting this to true for all components generates considerable load on the api server from watches.
		// Operations are additionally more complicated when all resources are watched.  Consider setting wait true for
		// relatively simple components, otherwise target specific resources with spec.healthChecks.
		wait: true | *false
		dependsOn: [for k, v in #DependsOn {v}]
	}
}

// #DependsOn stores all of the dependencies between components.  It's a struct to support merging across levels in the tree.
#DependsOn: {
	[Name=_]: {
		name: string | *"\(#InstancePrefix)-\(Name)"
	}
	...
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
	_namespace: string
	metadata: {
		name:      string | *"default"
		namespace: _namespace
	}
	spec: provider: {
		kubernetes: {
			remoteNamespace: _namespace
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

// #ClusterSpec is the specification of a holos platform cluster member.
#ClusterSpec: {
	// name is the cluster name.
	name: string
	// pool is the optional ceph pool of the cluster.
	pool?: string
	// region is the geographic region of the cluster.
	region?: string
	// primary is true if name matches the primaryCluster name
	primary: bool
}

// #Platform defines the primary lookup table for the platform.  Lookup keys should be limited to those defined in #KeyTags.
#Platform: {
	// org holds user defined values scoped organization wide.  A platform has one and only one organization.
	org: {
		// e.g. "example"
		name: string
		// e.g. "example.com"
		domain: string
		// e.g. "Example"
		displayName: string
		// e.g. "platform@example.com"
		contact: email: string
		// e.g. "platform@example.com"
		cloudflare: email: string
		// e.g. "example"
		github: orgs: primary: name: string
	}
	// Only one cluster may be primary at a time.  All others are standby.
	// Refer to [repo based standby](https://access.crunchydata.com/documentation/postgres-operator/latest/tutorials/backups-disaster-recovery/disaster-recovery#repo-based-standby)
	primaryCluster: {
		name: string
	}
	clusters: [Name=_]: #ClusterSpec & {
		name: string & Name
		if Name == primaryCluster.name {
			primary: true
		}
		if Name != primaryCluster.name {
			primary: false
		}
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

// ManagedNamespace is a namespace to manage across all clusters in the holos platform.
#ManagedNamespace: {
	namespace: {
		metadata: {
			name: string
			labels: [string]: string
		}
	}
	// clusterNames represents the set of clusters the namespace is managed on.  Usually all clusters.
	clusterNames: [...string]
}

// #ManagedNamepsaces is the union of all namespaces across all cluster types and optional services.
// Holos adopts the namespace sameness position of SIG Multicluster, refer to https://github.com/kubernetes/community/blob/dd4c8b704ef1c9c3bfd928c6fa9234276d61ad18/sig-multicluster/namespace-sameness-position-statement.md
#ManagedNamespaces: {
	[Name=_]: #ManagedNamespace & {
		namespace: metadata: name: Name
	}
}

// #Backups defines backup configuration.
// TODO: Consider the best place for this, possibly as part of the site platform config.  This represents the primary location for backups.
#Backups: {
	s3: {
		region:   string
		endpoint: string | *"s3.dualstack.\(region).amazonaws.com"
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
		ExternalSecret?: [Name=_]: #ExternalSecret & {_name: Name}
		VirtualService?: [Name=_]: #VirtualService & {metadata: name: Name}
		Issuer?: [Name=_]: #Issuer & {metadata: name: Name}
	}

	// apiObjectMap holds the marshalled representation of apiObjects
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

// #KubernetesObjects is the output schema of a single component.
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
	release: string | *name
	repository: {
		name?: string
		url?:  string
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
	// resources is the intermediate file name for api objects.
	resourcesFile: ResourcesFile
	// kustomizeFiles represents the files in a kustomize directory tree.
	kustomizeFiles: #KustomizeFiles.Files
	// enableHooks removes the --no-hooks flag from helm template
	enableHooks: true | *false
}

// #KustomizeBuild is a holos component that uses plain yaml files as the source of api objects for a holos component.
// Intended for upstream components like the CrunchyData Postgres Operator.  The holos cli is expected to execute kustomize build on the component directory to produce the rendered output.
#KustomizeBuild: {
	#OutputTypeMeta
	#APIObjects
	kind: "KustomizeBuild"
	metadata: name: #InstanceName
	// ksObjects holds the flux Kustomization objects for gitops.
	ksObjects: [...#Kustomization] | *[#Kustomization]
	// ksContent is the yaml representation of kustomization.
	ksContent: yaml.MarshalStream(ksObjects)
	// namespace defines the value passed to the helm --namespace flag
	namespace: #TargetNamespace
}

// #PlatformSpec is the output schema of a platform specification.
#PlatformSpec: {
	#OutputTypeMeta
	kind: "PlatformSpec"
}

// #SecretName is the name of a Secret, ususally coupling a Deployment to an ExternalSecret
#SecretName: string

// Cluster Domain is the cluster specific domain
#ClusterDomain: #InputKeys.cluster + "." + #Platform.org.domain

// #SidecarInject represents the istio sidecar inject label
#IstioSidecar: {
	"sidecar.istio.io/inject": "true"
	...
}

// #KustomizeTree represents a kustomize build.
#KustomizeFiles: {
	Objects: {
		"kustomization.yaml": #Kustomize
	}
	// Files holds the marshaled output holos writes to the filesystem
	Files: {
		for filename, obj in Objects {
			"\(filename)": yaml.Marshal(obj)
		}
		...
	}
}

// kustomization.yaml
#Kustomize: kc.#Kustomization & {
	apiVersion: "kustomize.config.k8s.io/v1beta1"
	kind:       "Kustomization"
	resources: [ResourcesFile]
	...
	if len(#KustomizePatches) > 0 {
		patches: [for v in #KustomizePatches {v}]
	}
}

#KustomizePatches: {
	[_]: #Patch
}

// #Patch is a kustomize patch
#Patch: kc.#Patch

// #DefaultSecurityContext is the holos default security context to comply with the restricted namespace policy.
// Refer to https://kubernetes.io/docs/concepts/security/pod-security-standards/#restricted
#DefaultSecurityContext: {
	securityContext: {
		allowPrivilegeEscalation: false
		runAsNonRoot:             true
		capabilities: drop: ["ALL"]
		seccompProfile: type: "RuntimeDefault"
	}
	...
}

// Certificate name should always match the secret name.
#Certificate: {
	metadata: name:   _
	spec: secretName: metadata.name
}

// #IsPrimaryCluster is true if the cluster being rendered is the primary cluster
// Used by the iam project to determine where https://login.example.com is active.
#IsPrimaryCluster: bool & #ClusterName == #Platform.primaryCluster.name

// By default, render kind: Skipped so holos knows to skip over intermediate cue files.
// This enables the use of holos render ./foo/bar/baz/... when bar contains intermediary constraints which are not complete components.
// Holos skips over these intermediary cue instances.
{} & #NoOutput

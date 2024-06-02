package holos

import (
	"encoding/yaml"
	v1 "github.com/holos-run/holos/api/v1alpha1"

	kc "sigs.k8s.io/kustomize/api/types"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	rbacv1 "k8s.io/api/rbac/v1"

	gwv1 "gateway.networking.k8s.io/gateway/v1"
	hrv1 "gateway.networking.k8s.io/httproute/v1"
	rgv1 "gateway.networking.k8s.io/referencegrant/v1beta1"

	ra "security.istio.io/requestauthentication/v1"
	ap "security.istio.io/authorizationpolicy/v1"

	is "cert-manager.io/issuer/v1"
	ci "cert-manager.io/clusterissuer/v1"
	certv1 "cert-manager.io/certificate/v1"

	ss "external-secrets.io/secretstore/v1beta1"
	es "external-secrets.io/externalsecret/v1beta1"

	pc "postgres-operator.crunchydata.com/postgrescluster/v1beta1"
)

// #Resources represents kubernetes api objects output along side a build plan.
// These resources are defined directly within CUE.
#Resources: {
	[Kind=string]: [NAME=string]: {
		kind: Kind
		metadata: name: string | *NAME
	}

	Namespace: [string]:             corev1.#Namespace
	ServiceAccount: [string]:        corev1.#ServiceAccount
	ConfigMap: [string]:             corev1.#ConfigMap
	Service: [string]:               corev1.#Service
	Deployment: [string]:            appsv1.#Deployment
	Job: [string]:                   batchv1.#Job
	CronJob: [string]:               batchv1.#CronJob
	ClusterRole: [string]:           rbacv1.#ClusterRole
	ClusterRoleBinding: [string]:    rbacv1.#ClusterRoleBinding
	Role: [string]:                  rbacv1.#Role
	RoleBinding: [string]:           rbacv1.#RoleBinding
	Issuer: [string]:                is.#Issuer
	ClusterIssuer: [string]:         ci.#ClusterIssuer
	Certificate: [string]:           certv1.#Certificate
	SecretStore: [string]:           ss.#SecretStore
	ExternalSecret: [string]:        es.#ExternalSecret
	HTTPRoute: [string]:             hrv1.#HTTPRoute
	ReferenceGrant: [string]:        rgv1.#ReferenceGrant
	PostgresCluster: [string]:       pc.#PostgresCluster
	RequestAuthentication: [string]: ra.#RequestAuthentication
	AuthorizationPolicy: [string]:   ap.#AuthorizationPolicy

	Gateway: [string]: gwv1.#Gateway & {
		spec: gatewayClassName: string | *"istio"
	}
}

#ReferenceGrant: rgv1.#ReferenceGrant & {
	spec: from: [{
		group:     "gateway.networking.k8s.io"
		kind:      "HTTPRoute"
		namespace: #IstioGatewaysNamespace
	}]
	spec: to: [{
		group: ""
		kind:  "Service"
	}]
}

// #Helm represents a holos build plan composed of one helm chart.
#Helm: {
	// Name represents the holos component name
	Name:      string
	Version:   string
	Namespace: string
	Resources: #Resources

	Repo: {
		name: string | *""
		url:  string | *""
	}

	Values: {...}

	Chart: v1.#HelmChart & {
		metadata: name: string | *Name
		namespace: string | *Namespace
		chart: name:       string | *Name
		chart: release:    chart.name
		chart: version:    string | *Version
		chart: repository: Repo

		// Render the values to yaml for holos to provide to helm.
		valuesContent: yaml.Marshal(Values)

		// Kustomize post-processor
		if EnableKustomizePostProcessor == true {
			// resourcesFile represents the file helm output is written two and
			// kustomize reads from.  Typically "resources.yaml" but referenced as a
			// constant to ensure the holos cli uses the same file.
			resourcesFile: v1.#ResourcesFile
			// kustomizeFiles represents the files in a kustomize directory tree.
			kustomizeFiles: v1.#FileContentMap
			for FileName, Object in KustomizeFiles {
				kustomizeFiles: "\(FileName)": yaml.Marshal(Object)
			}
		}

		apiObjectMap: (v1.#APIObjects & {apiObjects: Resources}).apiObjectMap
	}

	// EnableKustomizePostProcessor processes helm output with kustomize if true.
	EnableKustomizePostProcessor: true | *false
	// KustomizeFiles represents additional files to include in a Kustomization
	// resources list.  Useful to patch helm output.  The implementation is a
	// struct with filename keys and structs as values.  Holos encodes the struct
	// value to yaml then writes the result to the filename key.  Component
	// authors may then reference the filename in the kustomization.yaml resources
	// or patches lists.
	// Requires EnableKustomizePostProcessor: true.
	KustomizeFiles: {
		// Embed KustomizeResources
		KustomizeResources

		// The kustomization.yaml file must be included for kustomize to work.
		"kustomization.yaml": kc.#Kustomization & {
			apiVersion: "kustomize.config.k8s.io/v1beta1"
			kind:       "Kustomization"
			resources: [v1.#ResourcesFile, for FileName, _ in KustomizeResources {FileName}]
			patches: [for x in KustomizePatches {x}]
		}
	}
	// KustomizePatches represents patches to apply to the helm output.  Requires
	// EnableKustomizePostProcessor: true.
	KustomizePatches: [ArbitraryLabel=string]: kc.#Patch
	// KustomizeResources represents additional resources files to include in the
	// kustomize resources list.
	KustomizeResources: [FileName=string]: {...}

	// output represents the build plan provided to the holos cli.
	Output: v1.#BuildPlan & {
		spec: components: helmChartList: [Chart]
	}
}

// #Kustomize represents a holos build plan composed of one kustomize build.
#Kustomize: {
	// Name represents the holos component name
	Name: string

	Kustomization: v1.#KustomizeBuild & {
		metadata: name: string | *Name
	}

	// output represents the build plan provided to the holos cli.
	Output: v1.#BuildPlan & {
		spec: components: kustomizeBuildList: [Kustomization]
	}
}

// #Kubernetes represents a holos build plan composed of inline kubernetes api
// objects.
#Kubernetes: {
	// Name represents the holos component name
	Name:      string
	Namespace: string
	Resources: #Resources

	// output represents the build plan provided to the holos cli.
	Output: v1.#BuildPlan & {
		// resources is a map unlike other build plans which use a list.
		spec: components: resources: "\(Name)": {
			metadata: name: Name
			apiObjectMap: (v1.#APIObjects & {apiObjects: Resources}).apiObjectMap
		}
	}
}

package holos

import (
	"encoding/yaml"
	core "github.com/holos-run/holos/api/core/v1alpha2"

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

	app "argoproj.io/application/v1alpha1"

	cpv1 "pkg.crossplane.io/provider/v1"
	cpdrcv1beta1 "pkg.crossplane.io/deploymentruntimeconfig/v1beta1"
	cpfuncv1beta1 "pkg.crossplane.io/function/v1beta1"
	cpawspcv1beta1 "aws.upbound.io/providerconfig/v1beta1"
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

	// Crossplane resources
	DeploymentRuntimeConfig: [string]: cpdrcv1beta1.#DeploymentRuntimeConfig
	Provider: [string]:                cpv1.#Provider
	Function: [string]:                cpfuncv1beta1.#Function
	ProviderConfig: [string]:          cpawspcv1beta1.#ProviderConfig
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

	Chart: core.#HelmChart & {
		metadata: name:      string | *Name
		metadata: namespace: string | *Namespace
		chart: name:         string | *Name
		chart: release:      chart.name
		chart: version:      string | *Version
		chart: repository:   Repo

		// Render the values to yaml for holos to provide to helm.
		valuesContent: yaml.Marshal(Values)

		// Kustomize post-processor
		if EnableKustomizePostProcessor == true {
			// resourcesFile represents the file helm output is written two and
			// kustomize reads from.  Typically "resources.yaml" but referenced as a
			// constant to ensure the holos cli uses the same file.
			kustomize: resourcesFile: core.#ResourcesFile
			// kustomizeFiles represents the files in a kustomize directory tree.
			kustomize: kustomizeFiles: core.#FileContentMap
			for FileName, Object in KustomizeFiles {
				kustomize: kustomizeFiles: "\(FileName)": yaml.Marshal(Object)
			}
		}

		apiObjectMap: (#APIObjects & {apiObjects: Resources}).apiObjectMap
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
			resources: [core.#ResourcesFile, for FileName, _ in KustomizeResources {FileName}]
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
	Output: #BuildPlan & {
		_Name:      Name
		_Namespace: Namespace
		spec: components: helmChartList: [Chart]
	}
}

// #Kustomize represents a holos build plan composed of one kustomize build.
#Kustomize: {
	// Name represents the holos component name
	Name: string

	Kustomization: core.#KustomizeBuild & {
		metadata: name: string | *Name
	}

	// output represents the build plan provided to the holos cli.
	Output: #BuildPlan & {
		_Name: Name
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
	Output: #BuildPlan & {
		_Name:      Name
		_Namespace: Namespace
		// resources is a map unlike other build plans which use a list.
		spec: components: resources: "\(Name)": {
			metadata: name:      Name
			metadata: namespace: Namespace
			apiObjectMap: (#APIObjects & {apiObjects: Resources}).apiObjectMap
		}
	}
}

#BuildPlan: core.#BuildPlan & {
	_Name:       string
	_Namespace?: string
	let NAME = "gitops/\(_Name)"

	// Render the ArgoCD Application for GitOps.
	spec: components: resources: (NAME): {
		metadata: name: NAME
		if _Namespace != _|_ {
			metadata: namespace: _Namespace
		}

		deployFiles: (#Argo & {ComponentName: _Name}).deployFiles
	}
}

// #Argo represents an argocd Application resource for each component, written
// using the #HolosComponent.deployFiles field.
#Argo: {
	ComponentName: string

	Application: app.#Application & {
		metadata: name:      ComponentName
		metadata: namespace: "argocd"
		spec: {
			destination: server: "https://kubernetes.default.svc"
			project: "default"
			source: {
				path:           "\(_Platform.Model.argocd.deployRoot)/deploy/clusters/\(_ClusterName)/components/\(ComponentName)"
				repoURL:        _Platform.Model.argocd.repoURL
				targetRevision: _Platform.Model.argocd.targetRevision
			}
		}
	}

	// deployFiles represents the output files to write along side the component.
	deployFiles: "clusters/\(_ClusterName)/gitops/\(ComponentName).application.gen.yaml": yaml.Marshal(Application)
}

// #ArgoDefaultSyncPolicy represents the default argo sync policy.
#ArgoDefaultSyncPolicy: {
	automated: {
		prune:    bool | *true
		selfHeal: bool | *true
	}
	syncOptions: [
		"RespectIgnoreDifferences=true",
		"ServerSideApply=true",
	]
	retry: limit: number | *2
	retry: backoff: {
		duration:    string | *"5s"
		factor:      number | *2
		maxDuration: string | *"3m0s"
	}
}

// #APIObjects defines the output format for kubernetes api objects.  The holos
// cli expects the yaml representation of each api object in the apiObjectMap
// field.
#APIObjects: core.#APIObjects & {
	// apiObjects represents the un-marshalled form of each kubernetes api object
	// managed by a holos component.
	apiObjects: {
		[Kind=string]: {
			[string]: {
				kind: Kind
				...
			}
		}
		ConfigMap: [string]: corev1.#ConfigMap & {apiVersion: "v1"}
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
	}
}

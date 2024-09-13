package v1alpha3

import (
	"encoding/yaml"
	core "github.com/holos-run/holos/api/core/v1alpha3"
	kc "sigs.k8s.io/kustomize/api/types"

	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	batchv1 "k8s.io/api/batch/v1"

	app "argoproj.io/application/v1alpha1"
	ci "cert-manager.io/clusterissuer/v1"
	rgv1 "gateway.networking.k8s.io/referencegrant/v1beta1"
	certv1 "cert-manager.io/certificate/v1"
	hrv1 "gateway.networking.k8s.io/httproute/v1"
)

#Resources: {
	[Kind=string]: [InternalLabel=string]: {
		kind: Kind
		metadata: name: string | *InternalLabel
	}

	Certificate: [_]:        certv1.#Certificate
	ClusterIssuer: [_]:      ci.#ClusterIssuer
	ClusterRole: [_]:        rbacv1.#ClusterRole
	ClusterRoleBinding: [_]: rbacv1.#ClusterRoleBinding
	ConfigMap: [_]:          corev1.#ConfigMap
	CronJob: [_]:            batchv1.#CronJob
	Deployment: [_]:         appsv1.#Deployment
	HTTPRoute: [_]:          hrv1.#HTTPRoute
	Job: [_]:                batchv1.#Job
	Namespace: [_]:          corev1.#Namespace
	ReferenceGrant: [_]:     rgv1.#ReferenceGrant
	Role: [_]:               rbacv1.#Role
	RoleBinding: [_]:        rbacv1.#RoleBinding
	Service: [_]:            corev1.#Service
	ServiceAccount: [_]:     corev1.#ServiceAccount
	StatefulSet: [_]:        appsv1.#StatefulSet
	...
}

#Helm: {
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

	// ArgoConfig represents the ArgoCD GitOps integration for this Component.
	ArgoConfig: _

	// output represents the build plan provided to the holos cli.
	Output: #BuildPlan & {
		_Name:       Name
		_Namespace:  Namespace
		_ArgoConfig: ArgoConfig
		spec: components: helmChartList: [Chart]
	}
}

#BuildPlan: core.#BuildPlan & {
	_Name:       string
	_Namespace?: string
	_ArgoConfig: #ArgoConfig

	if _ArgoConfig.Enabled {
		let NAME = "gitops/\(_Name)"

		// Render the ArgoCD Application for GitOps as an additional Component of
		// the BuildPlan.
		spec: components: resources: (NAME): {
			metadata: name: NAME
			if _Namespace != _|_ {
				metadata: namespace: _Namespace
			}

			deployFiles: (#Argo & {ComponentName: _Name, ArgoConfig: _ArgoConfig}).deployFiles
		}
	}
}

// #Argo represents an argocd Application resource for each component, written
// using the #HolosComponent.deployFiles field.
#Argo: {
	ComponentName: string
	ArgoConfig:    #ArgoConfig

	Application: app.#Application & {
		metadata: name:      ComponentName
		metadata: namespace: "argocd"
		spec: {
			destination: server: "https://kubernetes.default.svc"
			project: "default"
			source: {
				path:           "\(ArgoConfig.DeployRoot)/deploy/clusters/\(ArgoConfig.ClusterName)/components/\(ComponentName)"
				repoURL:        ArgoConfig.RepoURL
				targetRevision: ArgoConfig.TargetRevision
			}
		}
	}

	// deployFiles represents the output files to write along side the component.
	deployFiles: "clusters/\(ArgoConfig.ClusterName)/gitops/\(ComponentName).application.gen.yaml": yaml.Marshal(Application)
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
	}

	// apiObjectMap holds the marshalled representation of apiObjects
	for kind, v in apiObjects {
		for name, obj in v {
			apiObjectMap: (kind): (name): yaml.Marshal(obj)
		}
	}
}

#Platform: {
	Name:  _
	Model: _
	Components: [string]: _
	Output: metadata: name: Name
	Output: spec: model:    Model
	Output: spec: components: [for c in Components {c}]
}

#Kustomize: {
	Name: _
	Kustomization: metadata: name: string | *Name
	Output: #BuildPlan & {
		_Name: Name
		spec: components: kustomizeBuildList: [Kustomization]
	}
}

#Kubernetes: {
	Name:      _
	Resources: #Resources

	Output: #BuildPlan & {
		_Name: Name
		// resources is a map unlike other build plans which use a list.
		spec: components: resources: (Name): {
			metadata: name: Name
			apiObjectMap: (#APIObjects & {apiObjects: Resources}).apiObjectMap
		}
	}
}

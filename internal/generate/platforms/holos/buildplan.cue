package holos

import (
	"encoding/yaml"
	v1 "github.com/holos-run/holos/api/v1alpha1"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"

	ci "cert-manager.io/clusterissuer/v1"
	certv1 "cert-manager.io/certificate/v1"

	ss "external-secrets.io/secretstore/v1beta1"
)

// #Resources represents kubernetes api objects output along side a build plan.
// These resources are defined directly within CUE.
#Resources: {
	[Kind=string]: [NAME=string]: {
		kind: Kind
		metadata: name: string | *NAME
	}

	Namespace: [string]:          corev1.#Namespace
	ServiceAccount: [string]:     corev1.#ServiceAccount
	ConfigMap: [string]:          corev1.#ConfigMap
	Job: [string]:                batchv1.#Job
	CronJob: [string]:            batchv1.#CronJob
	ClusterRole: [string]:        rbacv1.#ClusterRole
	ClusterRoleBinding: [string]: rbacv1.#ClusterRoleBinding
	Role: [string]:               rbacv1.#Role
	RoleBinding: [string]:        rbacv1.#RoleBinding
	ClusterIssuer: [string]:      ci.#ClusterIssuer
	Certificate: [string]:        certv1.#Certificate
	SecretStore: [string]:        ss.#SecretStore
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
		chart: version:    string | *Version
		chart: repository: Repo
		// Render the values to yaml for holos to provide to helm.
		valuesContent: yaml.Marshal(Values)

		apiObjectMap: (v1.#APIObjects & {apiObjects: Resources}).apiObjectMap
	}

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

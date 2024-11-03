package holos

import "github.com/holos-run/holos/api/author/v1alpha5:author"

#ComponentConfig: author.#ComponentConfig & {
	Name:      _Tags.component.name
	Path:      _Tags.component.path
	Resources: #Resources
}

// https://holos.run/docs/api/author/v1alpha5/#Kubernetes
#Kubernetes: close({
	#ComponentConfig
	author.#Kubernetes
})

// https://holos.run/docs/api/author/v1alpha5/#Kustomize
#Kustomize: close({
	#ComponentConfig
	author.#Kustomize
})

// https://holos.run/docs/api/author/v1alpha5/#Helm
#Helm: close({
	#ComponentConfig
	author.#Helm
})

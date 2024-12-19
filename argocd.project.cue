@if(!NoArgoCD)
package holos

import ap "argoproj.io/appproject/v1alpha1"

ArgoCD: #ArgoCD & {
	Version:   "2.13.2"
	Namespace: "argocd"
}

#ArgoCD: {
	Version:   string
	Namespace: string
}

// ArgoCD AppProject
#AppProject: ap.#AppProject & {
	metadata: name:      string
	metadata: namespace: string | *"argocd"
	spec: description:   string | *"Holos managed AppProject"
	spec: clusterResourceWhitelist: [{group: "*", kind: "*"}]
	spec: destinations: [{namespace: "*", server: "*"}]
	spec: sourceRepos: ["*"]
}

// Registration point for AppProjects
#AppProjects: [NAME=string]: #AppProject & {metadata: name: NAME}

// Register the ArgoCD Project namespaces and components
Projects: {
	argocd: {
		namespaces: (ArgoCD.Namespace): _
		components: {
			"app-projects": {
				name: "app-projects"
				path: "projects/argocd/components/app-projects"
			}
			"argocd-crds": {
				name: "argocd-crds"
				path: "projects/argocd/components/crds"
			}
			argocd: {
				name: "argocd"
				path: "projects/argocd/components/argocd"
			}
		}
	}
}

// Define at least the platform project.  Other components can register projects
// the same way from the root of the configuration.
AppProjects: #AppProjects

for PROJECT in Projects {
	AppProjects: (PROJECT.name): _
}

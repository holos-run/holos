package holos

import h "github.com/holos-run/holos/api/v1alpha1"

import "strings"

// #Projects is a map of all the projects in the platform.
#Projects: [Name=_]: #Project & {name: Name}

#Project: {
	name: string
	let ProjectName = name
	description: string
	environments: [Name=string]: #Environment & {
		name:    Name
		project: ProjectName
	}
	stages: [Name=string]: #Stage & {
		name:    Name
		project: ProjectName
	}

	// hosts are short hostnames to configure for the project.
	// Each value is routed to every environment in the project as a dns prefix.
	hosts: [Name=string]: #Host & {name: Name}
	// clusters are the cluster names the project is configured on.
	clusters: [Name=string]: #Cluster & {name: Name}
	clusterNames: [for c in clusters {c.name}]

	// managedNamespaces ensures project namespaces have SecretStores that can sync ExternalSecrets from the provisioner cluster.
	managedNamespaces: {
		// Define the shape of a managed namespace.
		[Name=_]: #ManagedNamespace & {
			namespace: metadata: name: Name
			clusterNames: ["provisioner", for c in clusters {c.name}]
		}

		// Manage a system namespace for each stage in the project.
		for stage in stages {
			for ns in stage.namespaces {
				(ns.name): _
			}
		}

		// Manage a namespace for each environment in the project.
		for env in environments {
			(env.namespace): _
		}
	}

	// features is YAGNI maybe? 
	features: [Name=string]: #Feature & {name: Name}
}

// #Cluster defines a cluster
#Cluster: name: string

// #Host defines a short hostname
#Host: name: string

#Environment: {
	// name uniquely identifies the environment within the scope of the project.
	name:      string
	project:   string
	stage:     string | "dev" | "prod"
	slug:      "\(name)-\(project)"
	namespace: "\(name)-\(project)"
	stageSlug: "\(stage)-\(project)"
	dnsSegments: [...string] | *[name]
}

#Stage: {
	name:    string
	project: string
	slug:    "\(name)-\(project)"
	// Manage a system namespace for each stage
	namespaces: [Name=_]: name: Name
	namespaces: "\(name)-\(project)-system": _
}

#Feature: {
	name:        string
	description: string
	enabled:     *true | false
}

#ProjectTemplate: {
	project: #Project

	// workload cluster resources
	workload: resources: [Name=_]: h.#KubernetesObjects & {
		metadata: name: Name
	}

	// provisioner cluster resources
	provisioner: resources: [Name=_]: h.#KubernetesObjects & {
		metadata: name: Name
	}
}

// #EnvHosts provides hostnames given a project and environment
#EnvHosts: {
	project: #Project & {name: env.project}
	env: #Environment

	hosts: {
		for host in project.hosts {
			let SEGMENTS = [host.name] + env.dnsSegments + [#Platform.org.domain]
			let HOST = strings.Join(SEGMENTS, ".")
			(HOST): {name: HOST}

			for cluster in project.clusters {
				let SEGMENTS = [host.name] + env.dnsSegments + [cluster.name, #Platform.org.domain]
				let HOST = strings.Join(SEGMENTS, ".")
				(HOST): {name: HOST}
			}
		}
	}
}

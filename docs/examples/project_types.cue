package holos

import h "github.com/holos-run/holos/api/v1alpha1"

import "strings"

// #Projects is a map of all the projects in the platform.
#Projects: [Name=_]: #Project & {name: Name}

_Projects: #Projects

// The platform project is required and where platform services reside.  ArgoCD, Grafana, Prometheus, etc...
#Projects: platform: _

#Project: {
	name: string
	// resourceId is the zitadel project Resource ID
	resourceId: number
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
	domain: string | *#Platform.org.domain

	// authProxyOrgDomain is the primary org domain for zitadel.
	authProxyOrgDomain: string | *#Platform.org.domain
	// authProxyIssuer is the issuer url
	authProxyIssuer: string | *"https://login.\(#Platform.org.domain)"

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
	features: authproxy: _
	features: httpbin:   _
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

	// envSegments are the env portion of the dns segments
	envSegments: [...string] | *[name]
	// stageSegments are the stage portion of the dns segments
	stageSegments: [...string] | *[stage]

	// #host provides a hostname
	// Refer to: https://github.com/holos-run/holos/issues/66#issuecomment-2027562626
	#host: {
		name:     string
		cluster?: string
		clusterSegments: [...string]
		if cluster != _|_ {
			clusterSegments: [cluster]
		}
		let SEGMENTS = envSegments + [name] + stageSegments + clusterSegments + [_Projects[project].domain]
		let NAMESEGMENTS = ["https"] + SEGMENTS
		host: {
			name: strings.Join(SEGMENTS, ".")
			port: {
				name:     strings.Replace(strings.Join(NAMESEGMENTS, "-"), ".", "-", -1)
				number:   443
				protocol: "HTTPS"
			}
		}
	}
}

#Stage: {
	name:    string
	project: string
	slug:    "\(name)-\(project)"
	// namespace is the system namespace for the project stage
	namespace: "\(name)-\(project)-system"
	// Manage a system namespace for each stage
	namespaces: [Name=_]: name: Name
	namespaces: (namespace): _
	// stageSegments are the stage portion of the dns segments
	stageSegments: [...string] | *[name]
	// authProxyClientID is the ClientID registered with the oidc issuer.
	authProxyClientID: string
	// extAuthzProviderName is the provider name in the mesh config
	extAuthzProviderName: "\(slug)-authproxy"
}

#Feature: {
	name:        string
	description: string
	enabled:     true | *false
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

// #EnvHosts provides hostnames given a project and environment.
// Refer to https://github.com/holos-run/holos/issues/66#issuecomment-2027562626
#EnvHosts: {
	project: #Project & {name: env.project}
	env: #Environment

	hosts: {
		for host in project.hosts {
			// globally scoped hostname
			let HOST = (env.#host & {name: host.name}).host
			(HOST.name): HOST

			// cluster scoped hostname
			for Cluster in project.clusters {
				let HOST = (env.#host & {name: host.name, cluster: Cluster.name}).host
				(HOST.name): HOST
			}
		}
	}
}

// #StageDomains provides hostnames given a project and stage.  Primarily for the
// auth proxy.
// Refer to https://github.com/holos-run/holos/issues/66#issuecomment-2027562626
#StageDomains: {
	// names are the leading prefix names to create hostnames for.
	// this is a two level list to support strings.Join()
	prefixes: [...[...string]] | *[[]]
	stage: #Stage
	project: #Project & {
		name: stage.project
	}

	// blank segment for the global domain plus each cluster in the project.
	let ClusterSegments = [[], for cluster in project.clusters {[cluster.name]}]

	hosts: {
		for prefix in prefixes {
			for ClusterSegment in ClusterSegments {
				let SEGMENTS = prefix + [project.name] + stage.stageSegments + ClusterSegment + [project.domain]
				let NAMESEGMENTS = ["https"] + SEGMENTS
				let HOSTNAME = strings.Join(SEGMENTS, ".")
				(HOSTNAME): {
					name: HOSTNAME
					port: {
						name:     strings.Replace(strings.Join(NAMESEGMENTS, "-"), ".", "-", -1)
						number:   443
						protocol: "HTTPS"
					}
				}
			}
		}
	}
}

// #StageDNSNames provides a CN and dnsNames for a project stage given a host
// Primarily for consolidating dns names into one certificate to avoid letsencrypt rate limits.
// Note, the project name is not part of the constructed dns names.
#StageCanonicalNames: {
	stage: #Stage
	project: #Project & {
		name: stage.project
	}

	// CanonicalNames is a mapping of CN to all dnsNames associated with the CN.
	CanonicalNames: {
		for host in project.hosts {
			let Name = host.name

			// Get a Canonical Name for the cert
			let SEGMENTS = [Name] + stage.stageSegments + [project.domain]
			let CN = strings.Join(SEGMENTS, ".")
			"\(CN)": {
				// CN is a member of it's own dnsNames entry.
				"\(CN)": CN

				// Map environment names and environment + cluster names.
				for Env in project.environments {
					if Env.stage == stage.name {
						let Host = (Env.#host & {name: Name}).host.name
						"\(Host)": Host

						for Cluster in project.clusters {
							let Host = (Env.#host & {name: Name, cluster: Cluster.name}).host.name
							"\(Host)": Host
						}
					}
				}
			}
		}
	}
}

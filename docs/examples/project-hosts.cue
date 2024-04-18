package holos

import "strings"

// #ProjectHosts represents all of the hosts associated with the project
// organized for use in Certificates, Gateways, VirtualServices.
#ProjectHosts: {
	project: #Project

	// Hosts map key fqdn to host to reduce into structs organized by stage,
	// canonical name, etc...  The flat nature and long list of properties is
	// intended to make it straight forward to derive another struct for Gateways,
	// VirtualServices, Certificates, AuthProxy cookie domains, etc...
	Hosts: {
		for Env in project.environments {
			for Host in project.hosts {
				// Global hostname, e.g. app.holos.run
				let CertInfo = (#MakeCertInfo & {
					host:       Host
					env:        Env
					domain:     project.domain
					clusterMap: project.clusters
				}).CertInfo

				"\(CertInfo.fqdn)": CertInfo

				// Cluster hostname, e.g. app.east1.holos.run, app.west1.holos.run
				for Cluster in project.clusters {
					let CertInfo = (#MakeCertInfo & {
						host:       Host
						env:        Env
						domain:     project.domain
						cluster:    Cluster.name
						clusterMap: project.clusters
					}).CertInfo

					"\(CertInfo.fqdn)": CertInfo
				}
			}
		}
	}
}

// #MakeCertInfo provides dns info for a certificate
// Refer to: https://github.com/holos-run/holos/issues/66#issuecomment-2027562626
#MakeCertInfo: {
	host:       #Host
	env:        #Environment
	domain:     string
	cluster:    string
	clusterMap: #ClusterMap

	let Stage = #StageInfo & {name: env.stage, project: env.project}
	let Env = env

	// DNS segments from left to right.
	let EnvSegments = env.envSegments

	WildcardSegments: [...string]
	if len(env.envSegments) > 0 {
		WildcardSegments: ["*"]
	}

	let HostSegments = [host.name]

	let StageSegments = env.stageSegments

	ClusterSegments: [...string]
	if cluster != _|_ {
		ClusterSegments: [cluster]
	}

	let DomainSegments = [domain]

	// Assemble the segments

	let FQDN = EnvSegments + HostSegments + StageSegments + ClusterSegments + DomainSegments
	let WILDCARD = WildcardSegments + HostSegments + StageSegments + ClusterSegments + DomainSegments
	let CANONICAL = HostSegments + StageSegments + DomainSegments

	CertInfo: #CertInfo & {
		fqdn:      strings.Join(FQDN, ".")
		wildcard:  strings.Join(WILDCARD, ".")
		canonical: strings.Join(CANONICAL, ".")

		project: name: Env.project
		stage: #StageOrEnvRef & {
			name:      Stage.name
			slug:      Stage.slug
			namespace: Stage.namespace
		}
		env: #StageOrEnvRef & {
			name:      Env.name
			slug:      Env.slug
			namespace: Env.namespace
		}

		if cluster != _|_ {
			// Host is valid on a single cluster.
			clusters: "\(cluster)": _
		}

		if cluster == _|_ {
			// Host is valid on all project clusters.
			clusters: clusterMap
		}
	}
}

// #CertInfo defines the attributes associated with a fully qualfied domain name
#CertInfo: {
	// fqdn is the fully qualified domain name, never a wildcard.
	fqdn: string
	// canonical is the canonical name this name may be an alternate name for.
	canonical: string
	// wildcard may replace the left most segment fqdn with a wildcard to consolidate cert dnsNames.  If not a wildcad, must be fqdn
	wildcard: string

	// Project, stage and env attributes for mapping and collecting.
	project: name: string

	stage: #StageOrEnvRef
	env:   #StageOrEnvRef

	// clusters represents the cluster names the fqdn is valid on.
	clusters: #ClusterMap
	// hosts are always valid on the provisioner cluster
	clusters: provisioner: _
}

#ClusterMap: [Name=string]: #Cluster & {name: Name}

#StageOrEnvRef: {
	name:      string
	slug:      string
	namespace: string
}

package holos

// import "encoding/yaml"

let Namespace = "prod-platform"

// FYI: kube-prometheus-stack is a large umbrella chart that brings in other large charts like
// [grafana](https://github.com/grafana/helm-charts/tree/main/charts/grafana).
// This may make affect maintainability.  Consider breaking the integration down into
// constituent charts represented as holos component instances.
spec: components: HelmChartList: [
	#HelmChart & {
		_dependsOn: "prod-secrets-stores": _
		namespace: Namespace

		metadata: name: "\(namespace)-obs"

		chart: {
			name:    "kube-prometheus-stack"
			release: "obs"
			version: "57.2.0"
			repository: {
				name: "prometheus-community"
				url:  "https://prometheus-community.github.io/helm-charts"
			}
		}

		_values: #KubePromStackValues & {
			fullnameOverride: "obs"

			// https://github.com/prometheus-community/helm-charts/tree/kube-prometheus-stack-57.2.0/charts/kube-prometheus-stack#alternatives
			prometheusOperator: admissionWebhooks: certManager: enabled: true

			prometheus: prometheusSpec: {
				// Pick up all ServiceMonitors in the monitoring namespace.  Note this means
				// only one Prometheus should be installed in the monitoring namespace.
				// See https://github.com/prometheus-community/helm-charts/issues/1911
				serviceMonitorSelectorNilUsesHelmValues: false

				// Unnecessary?
				// externalUrl: "https://\(Hosts.prometheus.name)"

				storageSpec: volumeClaimTemplate: spec: {
					accessModes: ["ReadWriteOnce"]
					resources: requests: storage: "10Gi"
				}
			}

			grafana: {
				enabled: true
				persistence: {
					enabled: true
					type:    "sts"
					accessModes: ["ReadWriteOnce"]
					size: "4Gi"
					finalizers: ["kubernetes.io/pvc-protection"]
				}

				"grafana.ini": {
					analytics: check_for_updates: false
					grafana_net: url:             "https://\(#Platform.org.emailDomain)"
					server: domain:               "\(Hosts.grafana.name)"
					// necessary to get the correct oidc redirect uri
					server: root_url:           "https://\(server.domain)/"
					auth: oauth_auto_login:     true
					auth: disable_login_form:   true
					auth: disable_signout_menu: true
					"auth.generic_oauth": {
						let OIDC = #Platform.oauthClients.grafana.spec
						enabled:             true
						name:                "Holos Platform"
						enable_login_token:  true
						auto_sign_up:        true
						client_id:           OIDC.clientID
						scopes:              OIDC.scopes
						auth_url:            OIDC.authorization_endpoint
						token_url:           OIDC.token_endpoint
						api_url:             OIDC.userinfo_endpoint
						use_pkce:            true
						name_attribute_path: "name"
						// TODO: Lift the admin, editor, and viewer group names up to the plaform config struct.
						role_attribute_path: "contains(groups[*], 'prod-cluster-admin') && 'Admin' || contains(groups[*], 'prod-cluster-editor') && 'Editor' || 'Viewer'"
					}
				}
			}
		}

		// Holos overlay objects
		apiObjectMap: OBJECTS.apiObjectMap
	},
]

let OBJECTS = #APIObjects & {
	apiObjects: {
		// ExternalSecret: "deploy-key": _
		VirtualService: (Hosts.prometheus.service): {
			metadata: name:      Hosts.prometheus.service
			metadata: namespace: Namespace
			spec: hosts: [for host in Hosts.prometheus.hosts {host.name}]
			spec: gateways: ["istio-ingress/default"]
			spec: http: [{route: [{destination: {
				host: "obs-prometheus.\(Namespace).svc.cluster.local"
				port: number: 9090
			}}]}]
		}
		VirtualService: (Hosts.grafana.service): {
			metadata: name:      Hosts.grafana.service
			metadata: namespace: Namespace
			spec: hosts: [for host in Hosts.grafana.hosts {host.name}]
			spec: gateways: ["istio-ingress/default"]
			spec: http: [{route: [{destination: {
				host: "obs-grafana.\(Namespace).svc.cluster.local"
				port: number: 80
			}}]}]
		}
	}
}

let Hosts = {
	[Service=_]: {
		name: string | *"\(service).\(#ClusterName).\(#Platform.org.domain)"
		let Name = name
		hosts: {
			cluster: name:  Name
			platform: name: "\(service).\(#Platform.org.domain)"
		}
		service: Service
	}

	prometheus: _
	grafana:    _
}

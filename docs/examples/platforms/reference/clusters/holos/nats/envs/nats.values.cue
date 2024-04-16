package holos

#NatsValues: {
	//###############################################################################
	// Global options
	//###############################################################################
	global: {
		image: {
			// global image pull policy to use for all container images in the chart
			// can be overridden by individual image pullPolicy
			pullPolicy: null
			// global list of secret names to use as image pull secrets for all pod specs in the chart
			// secrets must exist in the same namespace
			// https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/
			pullSecretNames: []
			// global registry to use for all container images in the chart
			// can be overridden by individual image registry
			registry: null
		}

		// global labels will be applied to all resources deployed by the chart
		labels: {}
	}

	//###############################################################################
	// Common options
	//###############################################################################
	// override name of the chart
	nameOverride: null
	// override full name of the chart+release
	fullnameOverride: null
	// override the namespace that resources are installed into
	namespaceOverride: null

	// reference a common CA Certificate or Bundle in all nats config `tls` blocks and nats-box contexts
	// note: `tls.verify` still must be set in the appropriate nats config `tls` blocks to require mTLS
	tlsCA: {
		enabled: false
		// set configMapName in order to mount an existing configMap to dir
		configMapName: null
		// set secretName in order to mount an existing secretName to dir
		secretName: null
		// directory to mount the configMap or secret to
		dir: "/etc/nats-ca-cert"
		// key in the configMap or secret that contains the CA Certificate or Bundle
		key: "ca.crt"
	}

	//###############################################################################
	// NATS Stateful Set and associated resources
	//###############################################################################
	//###########################################################
	// NATS config
	//###########################################################
	config: {
		cluster: {
			enabled: true | *false
			port:    6222
			// must be 2 or higher when jetstream is enabled
			replicas: 3

			// apply to generated route URLs that connect to other pods in the StatefulSet
			routeURLs: {
				// if both user and password are set, they will be added to route URLs
				// and the cluster authorization block
				user:     null
				password: null
				// set to true to use FQDN in route URLs
				useFQDN:          false
				k8sClusterDomain: "cluster.local"
			}

			tls: {
				enabled: true | *false
				// set secretName in order to mount an existing secret to dir
				secretName: null
				dir:        "/etc/nats-certs/cluster"
				cert:       "tls.crt"
				key:        "tls.key"
				// merge or patch the tls config
				// https://docs.nats.io/running-a-nats-service/configuration/securing_nats/tls
				merge: {}
				patch: []
			}

			// merge or patch the cluster config
			// https://docs.nats.io/running-a-nats-service/configuration/clustering/cluster_config
			merge: {}
			patch: []
		}

		jetstream: {
			enabled: true | *false

			fileStore: {
				enabled: true
				dir:     "/data"

				//###########################################################
				// stateful set -> volume claim templates -> jetstream pvc
				//###########################################################
				pvc: {
					enabled:          true
					size:             "10Gi"
					storageClassName: null

					// merge or patch the jetstream pvc
					// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#persistentvolumeclaim-v1-core
					merge: {}
					patch: []
					// defaults to "{{ include "nats.fullname" $ }}-js"
					name: null
				}

				// defaults to the PVC size
				maxSize: null
			}

			memoryStore: {
				enabled: false
				// ensure that container has a sufficient memory limit greater than maxSize
				maxSize: "1Gi"
			}

			// merge or patch the jetstream config
			// https://docs.nats.io/running-a-nats-service/configuration#jetstream
			merge: {}
			patch: []
		}

		nats: {
			port: 4222
			tls: {
				enabled: false
				// set secretName in order to mount an existing secret to dir
				secretName: null
				dir:        "/etc/nats-certs/nats"
				cert:       "tls.crt"
				key:        "tls.key"
				// merge or patch the tls config
				// https://docs.nats.io/running-a-nats-service/configuration/securing_nats/tls
				merge: {}
				patch: []
			}
		}

		leafnodes: {
			enabled: false
			port:    7422
			tls: {
				enabled: false
				// set secretName in order to mount an existing secret to dir
				secretName: null
				dir:        "/etc/nats-certs/leafnodes"
				cert:       "tls.crt"
				key:        "tls.key"
				// merge or patch the tls config
				// https://docs.nats.io/running-a-nats-service/configuration/securing_nats/tls
				merge: {}
				patch: []
			}

			// merge or patch the leafnodes config
			// https://docs.nats.io/running-a-nats-service/configuration/leafnodes/leafnode_conf
			merge: {}
			patch: []
		}

		websocket: {
			enabled: true | *false
			port:    8080
			tls: {
				enabled: false
				// set secretName in order to mount an existing secret to dir
				secretName: null
				dir:        "/etc/nats-certs/websocket"
				cert:       "tls.crt"
				key:        "tls.key"
				// merge or patch the tls config
				// https://docs.nats.io/running-a-nats-service/configuration/securing_nats/tls
				merge: {}
				patch: []
			}

			//###########################################################
			// ingress
			//###########################################################
			// service must be enabled also
			ingress: {
				enabled: false
				// must contain at least 1 host otherwise ingress will not be created
				hosts: []
				path:     "/"
				pathType: "Exact"
				// sets to the ingress class name
				className: null
				// set to an existing secret name to enable TLS on the ingress; applies to all hosts
				tlsSecretName: null

				// merge or patch the ingress
				// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#ingress-v1-networking-k8s-io
				merge: {}
				patch: []
				// defaults to "{{ include "nats.fullname" $ }}-ws"
				name: null
			}

			// merge or patch the websocket config
			// https://docs.nats.io/running-a-nats-service/configuration/websocket/websocket_conf
			merge: {}
			patch: []
		}

		mqtt: {
			enabled: false
			port:    1883
			tls: {
				enabled: false
				// set secretName in order to mount an existing secret to dir
				secretName: null
				dir:        "/etc/nats-certs/mqtt"
				cert:       "tls.crt"
				key:        "tls.key"
				// merge or patch the tls config
				// https://docs.nats.io/running-a-nats-service/configuration/securing_nats/tls
				merge: {}
				patch: []
			}

			// merge or patch the mqtt config
			// https://docs.nats.io/running-a-nats-service/configuration/mqtt/mqtt_config
			merge: {}
			patch: []
		}

		gateway: {
			enabled: false
			port:    7222
			tls: {
				enabled: false
				// set secretName in order to mount an existing secret to dir
				secretName: null
				dir:        "/etc/nats-certs/gateway"
				cert:       "tls.crt"
				key:        "tls.key"
				// merge or patch the tls config
				// https://docs.nats.io/running-a-nats-service/configuration/securing_nats/tls
				merge: {}
				patch: []
			}

			// merge or patch the gateway config
			// https://docs.nats.io/running-a-nats-service/configuration/gateways/gateway#gateway-configuration-block
			merge: {}
			patch: []
		}

		monitor: {
			enabled: true
			port:    8222
			tls: {
				// config.nats.tls must be enabled also
				// when enabled, monitoring port will use HTTPS with the options from config.nats.tls
				enabled: false
			}
		}

		profiling: {
			enabled: false
			port:    65432
		}

		resolver: {
			enabled: true | *false
			dir:     "/data/resolver"

			//###########################################################
			// stateful set -> volume claim templates -> resolver pvc
			//###########################################################
			pvc: {
				enabled:          true
				size:             "1Gi"
				storageClassName: null

				// merge or patch the pvc
				// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#persistentvolumeclaim-v1-core
				merge: {}
				patch: []
				// defaults to "{{ include "nats.fullname" $ }}-resolver"
				name: null
			}

			// merge or patch the resolver
			// https://docs.nats.io/running-a-nats-service/configuration/securing_nats/auth_intro/jwt/resolver
			merge: {
				type?:     string
				interval?: string
				timeout?:  string
			}
			patch: []
		}

		// adds a prefix to the server name, which defaults to the pod name
		// helpful for ensuring server name is unique in a super cluster
		serverNamePrefix: ""

		// merge or patch the nats config
		// https://docs.nats.io/running-a-nats-service/configuration
		// following special rules apply
		//  1. strings that start with << and end with >> will be unquoted
		//     use this for variables and numbers with units
		//  2. keys ending in $include will be switched to include directives
		//     keys are sorted alphabetically, use prefix before $includes to control includes ordering
		//     paths should be relative to /etc/nats-config/nats.conf
		// example:
		//
		//   merge:
		//     $include: ./my-config.conf
		//     zzz$include: ./my-config-last.conf
		//     server_name: nats
		//     authorization:
		//       token: << $TOKEN >>
		//     jetstream:
		//       max_memory_store: << 1GB >>
		//
		// will yield the config:
		// {
		//   include ./my-config.conf;
		//   "authorization": {
		//     "token": $TOKEN
		//   },
		//   "jetstream": {
		//     "max_memory_store": 1GB
		//   },
		//   "server_name": "nats",
		//   include ./my-config-last.conf;
		// }
		merge: {
			operator?:       string
			system_account?: string
			resolver_preload?: [string]: string
		}
		patch: []
	}

	//###########################################################
	// stateful set -> pod template -> nats container
	//###########################################################
	container: {
		image: {
			repository: "nats"
			tag:        "2.10.12-alpine"
			pullPolicy: null
			registry:   null
		}

		// container port options
		// must be enabled in the config section also
		// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#containerport-v1-core
		ports: {
			nats: {}
			leafnodes: {}
			websocket: {}
			mqtt: {}
			cluster: {}
			gateway: {}
			monitor: {}
			profiling: {}
		}

		// map with key as env var name, value can be string or map
		// example:
		//
		//   env:
		//     GOMEMLIMIT: 7GiB
		//     TOKEN:
		//       valueFrom:
		//         secretKeyRef:
		//           name: nats-auth
		//           key: token
		env: {}

		// merge or patch the container
		// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#container-v1-core
		merge: {}
		patch: []
	}

	//###########################################################
	// stateful set -> pod template -> reloader container
	//###########################################################
	reloader: {
		enabled: true
		image: {
			repository: "natsio/nats-server-config-reloader"
			tag:        "0.14.1"
			pullPolicy: null
			registry:   null
		}

		// env var map, see nats.env for an example
		env: {}

		// all nats container volume mounts with the following prefixes
		// will be mounted into the reloader container
		natsVolumeMountPrefixes: ["/etc/"]

		// merge or patch the container
		// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#container-v1-core
		merge: {}
		patch: []
	}

	//###########################################################
	// stateful set -> pod template -> prom-exporter container
	//###########################################################
	// config.monitor must be enabled
	promExporter: {
		enabled: true | *false
		image: {
			repository: "natsio/prometheus-nats-exporter"
			tag:        "0.14.0"
			pullPolicy: null
			registry:   null
		}

		port: 7777
		// env var map, see nats.env for an example
		env: {}

		// merge or patch the container
		// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#container-v1-core
		merge: {}
		patch: []

		//###########################################################
		// prometheus pod monitor
		//###########################################################
		podMonitor: {
			enabled: true | *false

			// merge or patch the pod monitor
			// https://prometheus-operator.dev/docs/operator/api/#monitoring.coreos.com/v1.PodMonitor
			merge: {}
			patch: []
			// defaults to "{{ include "nats.fullname" $ }}"
			name: null
		}
	}

	//###########################################################
	// service
	//###########################################################
	service: {
		enabled: true

		// service port options
		// additional boolean field enable to control whether port is exposed in the service
		// must be enabled in the config section also
		// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#serviceport-v1-core
		ports: {
			nats: enabled:      true
			leafnodes: enabled: true
			websocket: enabled: true
			mqtt: enabled:      true
			cluster: enabled:   false
			gateway: enabled:   false
			monitor: enabled:   false
			profiling: enabled: false
		}

		// merge or patch the service
		// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#service-v1-core
		merge: {}
		patch: []
		// defaults to "{{ include "nats.fullname" $ }}"
		name: null
	}

	//###########################################################
	// other nats extension points
	//###########################################################
	// stateful set
	statefulSet: {
		// merge or patch the stateful set
		// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#statefulset-v1-apps
		merge: {}
		patch: []
		// defaults to "{{ include "nats.fullname" $ }}"
		name: null
	}

	// stateful set -> pod template
	podTemplate: {
		// adds a hash of the ConfigMap as a pod annotation
		// this will cause the StatefulSet to roll when the ConfigMap is updated
		configChecksumAnnotation: true

		// map of topologyKey: topologySpreadConstraint
		// labelSelector will be added to match StatefulSet pods
		//
		// topologySpreadConstraints:
		//   kubernetes.io/hostname:
		//     maxSkew: 1
		//
		topologySpreadConstraints: {}

		// merge or patch the pod template
		// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#pod-v1-core
		merge: {}
		patch: []
	}

	// headless service
	headlessService: {
		// merge or patch the headless service
		// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#service-v1-core
		merge: {}
		patch: []
		// defaults to "{{ include "nats.fullname" $ }}-headless"
		name: null
	}

	// config map
	configMap: {
		// merge or patch the config map
		// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#configmap-v1-core
		merge: {}
		patch: []
		// defaults to "{{ include "nats.fullname" $ }}-config"
		name: null
	}

	// pod disruption budget
	podDisruptionBudget: {
		enabled: true
		// merge or patch the pod disruption budget
		// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#poddisruptionbudget-v1-policy
		merge: {}
		patch: []
		// defaults to "{{ include "nats.fullname" $ }}"
		name: null
	}

	// service account
	serviceAccount: {
		enabled: false
		// merge or patch the service account
		// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#serviceaccount-v1-core
		merge: {}
		patch: []
		// defaults to "{{ include "nats.fullname" $ }}"
		name: null
	}

	//###########################################################
	// natsBox
	//
	// NATS Box Deployment and associated resources
	//###########################################################
	natsBox: {
		enabled: true

		//###########################################################
		// NATS contexts
		//###########################################################
		contexts: {
			default: {
				creds: {
					// set contents in order to create a secret with the creds file contents
					contents: null
					// set secretName in order to mount an existing secret to dir
					secretName: null
					// defaults to /etc/nats-creds/<context-name>
					dir: null
					key: "nats.creds"
				}
				nkey: {
					// set contents in order to create a secret with the nkey file contents
					contents: null
					// set secretName in order to mount an existing secret to dir
					secretName: null
					// defaults to /etc/nats-nkeys/<context-name>
					dir: null
					key: "nats.nk"
				}
				// used to connect with client certificates
				tls: {
					// set secretName in order to mount an existing secret to dir
					secretName: null
					// defaults to /etc/nats-certs/<context-name>
					dir:  null
					cert: "tls.crt"
					key:  "tls.key"
				}

				// merge or patch the context
				// https://docs.nats.io/using-nats/nats-tools/nats_cli#nats-contexts
				merge: {}
				patch: []
			}
		}

		// name of context to select by default
		defaultContextName: "default"

		//###########################################################
		// deployment -> pod template -> nats-box container
		//###########################################################
		container: {
			image: {
				repository: "natsio/nats-box"
				tag:        "0.14.2"
				pullPolicy: null
				registry:   null
			}

			// env var map, see nats.env for an example
			env: {}

			// merge or patch the container
			// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#container-v1-core
			merge: {}
			patch: []
		}

		//###########################################################
		// other nats-box extension points
		//###########################################################
		// deployment
		deployment: {
			// merge or patch the deployment
			// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#deployment-v1-apps
			merge: {}
			patch: []
			// defaults to "{{ include "nats.fullname" $ }}-box"
			name: null
		}

		// deployment -> pod template
		podTemplate: {
			// merge or patch the pod template
			// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#pod-v1-core
			merge: {}
			patch: []
		}

		// contexts secret
		contextsSecret: {
			// merge or patch the context secret
			// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#secret-v1-core
			merge: {}
			patch: []
			// defaults to "{{ include "nats.fullname" $ }}-box-contexts"
			name: null
		}

		// contents secret
		contentsSecret: {
			// merge or patch the contents secret
			// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#secret-v1-core
			merge: {}
			patch: []
			// defaults to "{{ include "nats.fullname" $ }}-box-contents"
			name: null
		}

		// service account
		serviceAccount: {
			enabled: false
			// merge or patch the service account
			// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#serviceaccount-v1-core
			merge: {}
			patch: []
			// defaults to "{{ include "nats.fullname" $ }}-box"
			name: null
		}
	}

	//###############################################################################
	// Extra user-defined resources
	//###############################################################################
	//
	// add arbitrary user-generated resources
	// example:
	//
	// config:
	//   websocket:
	//     enabled: true
	// extraResources:
	// - apiVersion: networking.istio.io/v1beta1
	//   kind: VirtualService
	//   metadata:
	//     name:
	//       $tplYaml: >
	//         {{ include "nats.fullname" $ | quote }}
	//     labels:
	//       $tplYaml: |
	//         {{ include "nats.labels" $ }}
	//   spec:
	//     hosts:
	//     - demo.nats.io
	//     gateways:
	//     - my-gateway
	//     http:
	//     - name: default
	//       match:
	//       - name: root
	//         uri:
	//           exact: /
	//       route:
	//       - destination:
	//           host:
	//             $tplYaml: >
	//               {{ .Values.service.name | quote }}
	//           port:
	//             number:
	//               $tplYaml: >
	//                 {{ .Values.config.websocket.port }}
	//
	extraResources: []
}

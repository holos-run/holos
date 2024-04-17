package holos

// for Project in _Projects {
// 	spec: components: resources: (#ProjectTemplate & {project: Project}).workload.resources
// }

let Namespace = "jeff-holos"

#Kustomization: spec: targetNamespace: Namespace

spec: components: HelmChartList: [
	#HelmChart & {
		metadata: name: "jeff-holos-nats"
		namespace: Namespace
		_dependsOn: "prod-secrets-stores": _
		chart: {
			name:       "nats"
			version:    "1.1.10"
			repository: NatsRepository
		}
		_values: #NatsValues & {
			config: {
				// https://github.com/nats-io/k8s/tree/main/helm/charts/nats#operator-mode-with-nats-resolver
				resolver: enabled: true
				resolver: merge: {
					type:     "full"
					interval: "2m"
					timeout:  "1.9s"
				}
				merge: {
					operator:       "eyJ0eXAiOiJKV1QiLCJhbGciOiJlZDI1NTE5LW5rZXkifQ.eyJqdGkiOiJUSElBTDM2NUtOS0lVVVJDMzNLNFJGQkJVRlFBSTRLS0NQTDJGVDZYVjdNQVhWU1dFNElRIiwiaWF0IjoxNzEzMjIxMzE1LCJpc3MiOiJPREtQM0RZTzc3T1NBRU5IU0FFR0s3WUNFTFBYT1FFWUI3RVFSTVBLWlBNQUxINE5BRUVLSjZDRyIsIm5hbWUiOiJIb2xvcyIsInN1YiI6Ik9ES1AzRFlPNzdPU0FFTkhTQUVHSzdZQ0VMUFhPUUVZQjdFUVJNUEtaUE1BTEg0TkFFRUtKNkNHIiwibmF0cyI6eyJ0eXBlIjoib3BlcmF0b3IiLCJ2ZXJzaW9uIjoyfX0.dQURTb-zIQMc-OYd9328oY887AEnvog6gOXY1-VCsDG3L89nq5x_ks4ME7dJ4Pn-Pvm2eyBi1Jx6ubgkthHgCQ"
					system_account: "ADIQCYK4K3OKTPODGCLI4PDQ6XBO52MISBPTAIDESEJMLZCMNULDKCCY"
					resolver_preload: {
						// NOTEL: Make sure you do not include the trailing , in the SYS_ACCOUNT_JWT
						"ADIQCYK4K3OKTPODGCLI4PDQ6XBO52MISBPTAIDESEJMLZCMNULDKCCY": "eyJ0eXAiOiJKV1QiLCJhbGciOiJlZDI1NTE5LW5rZXkifQ.eyJqdGkiOiI2SEVMNlhKSUdWUElMNFBURVI1MkUzTkFITjZLWkVUUUdFTlFVS0JWRzNUWlNLRzVLT09RIiwiaWF0IjoxNzEzMjIxMzE1LCJpc3MiOiJPREtQM0RZTzc3T1NBRU5IU0FFR0s3WUNFTFBYT1FFWUI3RVFSTVBLWlBNQUxINE5BRUVLSjZDRyIsIm5hbWUiOiJTWVMiLCJzdWIiOiJBRElRQ1lLNEszT0tUUE9ER0NMSTRQRFE2WEJPNTJNSVNCUFRBSURFU0VKTUxaQ01OVUxES0NDWSIsIm5hdHMiOnsibGltaXRzIjp7InN1YnMiOi0xLCJkYXRhIjotMSwicGF5bG9hZCI6LTEsImltcG9ydHMiOi0xLCJleHBvcnRzIjotMSwid2lsZGNhcmRzIjp0cnVlLCJjb25uIjotMSwibGVhZiI6LTF9LCJkZWZhdWx0X3Blcm1pc3Npb25zIjp7InB1YiI6e30sInN1YiI6e319LCJhdXRob3JpemF0aW9uIjp7fSwidHlwZSI6ImFjY291bnQiLCJ2ZXJzaW9uIjoyfX0.TiGIk8XON394D9SBEowGHY_nTeOyHiM-ihyw6HZs8AngOnYPFXH9OVjsaAf8Poa2k_V84VtH7yVNgNdjBgduDA"
					}
				}
				cluster: enabled:   true
				jetstream: enabled: true
				websocket: enabled: true
				monitor: enabled:   true
			}
			promExporter: enabled: true
			promExporter: podMonitor: enabled: true
		}
	},
	#HelmChart & {
		metadata: name: "jeff-holos-nack"
		namespace: Namespace
		_dependsOn: "jeff-holos-nats": _
		chart: {
			name:       "nack"
			version:    "0.25.2"
			repository: NatsRepository
		}
	},
]

let NatsRepository = {
	name: "nats"
	url:  "https://nats-io.github.io/k8s/helm/charts/"
}

package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

let Objects = {
	Name:      "{{ .Name }}"
	Namespace: #Istio.Gateway.Namespace

	Resources: {
		// The default gateway with all listeners attached to tls certs.
		Gateway: default: {
			metadata: namespace: Namespace

			let Listeners = {
				http: {
					name:     "http"
					protocol: "HTTP"
					port:     80
					allowedRoutes: namespaces: from: "Same"
				}
				https: {
					name:     "https"
					protocol: "HTTPS"
					port:     443
					allowedRoutes: namespaces: from: "Same"
					tls: mode: "Terminate"
					tls: certificateRefs: [{
						kind: "Secret"
						name: "gateway-cert"
					}]
				}
			}

			spec: listeners: [for x in Listeners {x}]
		}

		// Manage a simple cert for example.com and *.example.com
		Certificate: "gateway-cert": {
			metadata: name:      "gateway-cert"
			metadata: namespace: Namespace
			spec: commonName:    #Platform.Domain
			spec: dnsNames: [spec.commonName, "*.\(spec.commonName)"]
			spec: secretName: metadata.name
			spec: issuerRef: {
				kind: "ClusterIssuer"
				name: "local-ca"
			}
		}

		// Manage a service account to prevent ArgoCD from pruning it.
		ServiceAccount: "default-istio": {
			metadata: namespace: Namespace
			metadata: labels: {
				"gateway.istio.io/managed":               "istio.io-gateway-controller"
				"gateway.networking.k8s.io/gateway-name": "default"
				"istio.io/gateway-name":                  "default"
			}
		}
	}
}

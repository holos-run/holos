package holos

holos: {
  kind: "BuildPlan",
  apiVersion: "v1alpha5",
  metadata: name: "basic",
  spec: {
    artifacts: [
      {
        artifact: "components/\(metadata.name)/resources.gen.yaml"
        generators: [
          {
            "kind": "Resources",
            "output": artifact,
            "resources": {
              "Deployment": {
                "httpbin": {
                  "apiVersion": "apps/v1",
                  "kind": "Deployment",
                  "metadata": {
                    "name": "httpbin",
                    "namespace": "httpbin"
                  },
                  "spec": {
                    "replicas": 1,
                    "selector": {
                      "matchLabels": {
                        "app.kubernetes.io/name": "httpbin"
                      }
                    },
                    "template": {
                      "metadata": {
                        "labels": {
                          "app.kubernetes.io/name": "httpbin"
                        }
                      },
                      "spec": {
                        "containers": [
                          {
                            "image": "quay.io/holos/mccutchen/go-httpbin",
                            "livenessProbe": {
                              "httpGet": {
                                "path": "/status/200",
                                "port": "http"
                              }
                            },
                            "name": "httpbin",
                            "ports": [
                              {
                                "containerPort": 8080,
                                "name": "http",
                                "protocol": "TCP"
                              }
                            ],
                            "readinessProbe": {
                              "httpGet": {
                                "path": "/status/200",
                                "port": "http"
                              }
                            },
                          }
                        ]
                      }
                    }
                  }
                }
              },
            },
          }
        ],
      }
    ]
  }
}

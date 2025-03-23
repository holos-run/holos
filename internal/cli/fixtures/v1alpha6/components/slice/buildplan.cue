package holos

// Example to exercise and develop the slice transformer as a concrete use case
// for a follow up arbitrary command transformer.

// holos show buildplans --format=json | pbcopy
holos: {
  "kind": "BuildPlan",
  "apiVersion": "v1alpha6",
  "metadata": {
    "name": "slice",
    "labels": {
      "holos.run/component.name": "slice"
    },
    "annotations": {
      "app.holos.run/description": "slice transformer"
    }
  },
  "spec": {
    "artifacts": [
      {
        "generators": [
          {
            "kind": "Resources",
            "output": "resources.gen.yaml",
            "resources": {
              "Deployment": {
                "httpbin": {
                  "apiVersion": "apps/v1",
                  "kind": "Deployment",
                  "metadata": {
                    "name": "httpbin",
                    "namespace": "httpbin-demo"
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
                            "resources": {}
                          }
                        ]
                      }
                    }
                  }
                }
              },
              "Service": {
                "httpbin": {
                  "apiVersion": "v1",
                  "kind": "Service",
                  "metadata": {
                    "name": "httpbin",
                    "namespace": "httpbin-demo"
                  },
                  "spec": {
                    "ports": [
                      {
                        "appProtocol": "http",
                        "name": "http",
                        "port": 80,
                        "protocol": "TCP",
                        "targetPort": "http"
                      }
                    ],
                    "selector": {
                      "app.kubernetes.io/name": "httpbin"
                    }
                  }
                }
              }
            },
          }
        ],
      }
    ]
  }
}

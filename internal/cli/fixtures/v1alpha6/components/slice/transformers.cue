package holos

holos: {
  buildContext: _
  "spec": {
    "artifacts": [
      {
        artifact: "components/slice",
        "transformers": [
          {
            "kind": "Kustomize",
            "inputs": [
              "resources.gen.yaml"
            ],
            "output": "slice.gen.yaml",
            "kustomize": {
              "kustomization": {
                "apiVersion": "kustomize.config.k8s.io/v1beta1",
                "kind": "Kustomization",
                "resources": [
                  "resources.gen.yaml"
                ]
              }
            }
          },
          {
            "kind": "Command"
            "inputs": ["slice.gen.yaml"]
            "output": artifact
            "command": {
              "args": [
                "kubectl-slice",
                "-f",
                "\(buildContext.tempDir)/slice.gen.yaml",
                "-o",
                "\(buildContext.tempDir)/\(artifact)",
              ]
            }
          }
        ]
      }
    ]
  }
}

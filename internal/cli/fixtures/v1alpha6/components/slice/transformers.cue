package holos

// Focus in on the slice transformer
let ARTIFACT = "components/slice"

holos: {
  context: _
  "spec": {
    "artifacts": [
      {
        artifact: ARTIFACT,
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
            "output": ARTIFACT
            "command": {
              "args": [
                "kubectl-slice",
                "-f",
                "\(context.tempDir)/slice.gen.yaml",
                "-o",
                "\(context.tempDir)/\(ARTIFACT)",
              ]
            }
          }
        ]
      }
    ]
  }
}

package holos

holos: {
  "kind": "BuildPlan",
  "apiVersion": "v1alpha5",
  "metadata": {
    "name": "example"
  },
  "spec": {
    "artifacts": [
      {
        "artifact": "v1alpha5/example/example.gen.yaml",
        "generators": [
          {
            "kind": "Resources",
            "output": "v1alpha5/example/example.gen.yaml"
          }
        ],
      }
    ]
  }
}

kind: BuildPlan
apiVersion: v1alpha5
metadata:
  name: no-name
spec:
  artifacts:
    - artifact: components/no-name/no-name.gen.yaml
      generators:
        - kind: Resources
          output: resources.gen.yaml
          resources: {}
        - kind: File
          output: httpbin.yaml
          file:
            source: httpbin.yaml
      transformers:
        - kind: Kustomize
          inputs:
            - resources.gen.yaml
            - httpbin.yaml
          output: components/no-name/no-name.gen.yaml
          kustomize:
            kustomization:
              apiVersion: kustomize.config.k8s.io/v1beta1
              kind: Kustomization
              labels:
                - includeSelectors: false
                  pairs:
                    app.kubernetes.io/name: httpbin
              patches: []
              images:
                - name: mccutchen/go-httpbin
              resources:
                - resources.gen.yaml
                - httpbin.yaml
      validators: []

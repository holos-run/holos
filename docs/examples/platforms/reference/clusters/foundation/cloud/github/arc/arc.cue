package holos

// GitHub Actions Runner Controller
#InputKeys: project: "github"
#DependsOn: Namespaces: name: "prod-secrets-namespaces"

#ARCSystemNamespace: "arc-system"
#HelmChart: namespace: #TargetNamespace
#HelmChart: chart: version: "0.8.3"

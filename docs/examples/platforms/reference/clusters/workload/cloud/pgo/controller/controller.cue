package holos

#DependsOn: Namespaces: name: "prod-secrets-namespaces"
#DependsOn: CRDS: name: "\(#InstancePrefix)-crds"
#InputKeys: component: "controller"
{} & #KustomizeBuild

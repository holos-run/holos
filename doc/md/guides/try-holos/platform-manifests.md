# Platform Manifests

This document provides an example of how Holos uses CUE and Helm to unify and
render the platform configuration.  It refers to the manifests rendered in the
[Try Holos Locally](/docs/guides/try-holos/) guide.

Take a moment to review the manifests `holos` rendered to build the platform.

### ArgoCD Application

Note the Git URL in the Platform Model is used to derive the ArgoCD
`Application` resource for all of the platform components.

```yaml
# deploy/clusters/workload/gitops/namespaces.application.gen.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: namespaces
  namespace: argocd
spec:
  destination:
    server: https://kubernetes.default.svc
  project: default
  source:
    # highlight-next-line
    path: /deploy/clusters/workload/components/namespaces
    # highlight-next-line
    repoURL: https://github.com/holos-run/holos-k3d.git
    # highlight-next-line
    targetRevision: HEAD
```

One ArgoCD `Application` resource is produced for each Holos component by
default.  The CUE definition which produces the rendered output is defined in
`buildplan.cue` around line 222.

:::tip

Note how CUE does not use error-prone text templates, the language is well
specified and typed which reduces errors when unifying the configuration with
the Platform Model in the following `#Argo` definition.

:::

```cue
// buildplan.cue

// #Argo represents an argocd Application resource for each component, written
// using the #HolosComponent.deployFiles field.
#Argo: {
	ComponentName: string

	Application: app.#Application & {
		metadata: name:      ComponentName
		metadata: namespace: "argocd"
		spec: {
			destination: server: "https://kubernetes.default.svc"
			project: "default"
			source: {
        // highlight-next-line
				path:           "\(_Platform.Model.argocd.deployRoot)/deploy/clusters/\(_ClusterName)/components/\(ComponentName)"
        // highlight-next-line
				repoURL:        _Platform.Model.argocd.repoURL
        // highlight-next-line
				targetRevision: _Platform.Model.argocd.targetRevision
			}
		}
	}

	// deployFiles represents the output files to write along side the component.
	deployFiles: "clusters/\(_ClusterName)/gitops/\(ComponentName).application.gen.yaml": yaml.Marshal(Application)
}
```

### Helm Chart

The `cert-manger` component renders using the upstream Helm chart.  The build
plan that defines the helm chart to use along with the values to provide looks
like the following.

:::tip

Holos fully supports your existing Helm charts.  Consider leveraging `holos` as
an alternative to umbrella charts.

:::

```cue
// components/cert-manager/cert-manager.cue
package holos

// Produce a helm chart build plan.
(#Helm & Chart).Output

let Chart = {
	Name:      "cert-manager"
	Version:   "1.14.5"
	Namespace: "cert-manager"

	Repo: name: "jetstack"
	Repo: url:  "https://charts.jetstack.io"

  // highlight-next-line
	Values: {
		installCRDs: true
		startupapicheck: enabled: false
		// Must not use kube-system on gke autopilot.  GKE Warden blocks access.
    // highlight-next-line
		global: leaderElection: namespace: Namespace

		// https://cloud.google.com/kubernetes-engine/docs/concepts/autopilot-resource-requests#min-max-requests
		resources: requests: {
			cpu:                 "250m"
			memory:              "512Mi"
			"ephemeral-storage": "100Mi"
		}
    // highlight-next-line
		webhook: resources:        Values.resources
    // highlight-next-line
		cainjector: resources:     Values.resources
    // highlight-next-line
		startupapicheck: resource: Values.resources

		// https://cloud.google.com/kubernetes-engine/docs/how-to/autopilot-spot-pods
		nodeSelector: {
			"kubernetes.io/os": "linux"
			if _ClusterName == "management" {
				"cloud.google.com/gke-spot": "true"
			}
		}
		webhook: nodeSelector:         Values.nodeSelector
		cainjector: nodeSelector:      Values.nodeSelector
		startupapicheck: nodeSelector: Values.nodeSelector
	}
}
```

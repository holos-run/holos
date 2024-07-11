# Try Holos with k3d

Learn how to configure and deploy the Holos reference platform to your local host with k3d.

---

This guide assumes commands are run from your local host.  Capitalized terms
have specific definitions described in the [Glossary](/docs/glossary)

## Outcome

At the end of this guide you'll have built a development platform which
integrates the following components together to implement Zero Trust security:

 1. ArgoCD to review and apply platform configuration changes.
 2. Istio service mesh with mTLS encryption.
 3. ZITADEL to provide single sign-on identity tokens with multi factor authentication.

The platform running on your local host will configure Istio to authenticate and
authorize requests using an oidc id token issued by ZITADEL _before_ the request
ever reaches ArgoCD.

:::tip

With Holos, developers don't need to write authentication or authorization logic
for many use cases.  Single sign-on and role based access control are provided
by the platform itself for all service running in the platform using
standardized policies.

:::

The `k3d` platform is intentionally less holistic than the Holos reference
platform.  Larger, more holistic integrations are traded in for a shorter and
smoother on-ramp to learn about some of the value of Holos:

 1. Holos wraps unmodified Helm charts provided by software vendors.
 2. Holos eliminates the need to template yaml.
 3. Holos is composeable, scaling down to local host and up to multi-cloud and multi-cluster.
 4. The Zero Trust security model implemented by the reference platform.
 5. Configuration unification with CUE.

## Requirements

You'll need the following tools installed on your local host to complete this guide.

 1. [k3d](https://k3d.io/#installation)
 2. [Docker](https://docs.docker.com/get-docker/) to use k3d.
 2. [holos](/docs/tutorial/install) to build the platform configuration.
 3. [kubectl](https://kubernetes.io/docs/tasks/tools/) to interact with the Kubernetes cluster.
 4. [helm](https://helm.sh/docs/intro/install/) to render Holos components that integrate vendor provided Helm charts.

## Register with Holos

Register an account with the Holos web service, which is necessary to save
platform configuration values using a simple web form.

```bash
holos register user
```

## Create the Platform

Create the platform in the Holos web service to store the Platform Form and the form
values which represents the Platform Model.

```bash
holos create platform --name k3d --display-name "Try Holos Locally"
```

## Generate the Platform

Holos builds the platform by building each component of the platform into fully
rendered Kubernetes configuration resources.  Generate the source code for the
platform in a blank local directory.  This directory is named `holos-infra` by
convention because it represents the Holos managed platform infrastructure.

Create a new Git repository to store the platform code:

```bash
mkdir holos-k3d
cd holos-k3d
git init .
```

Generate the platform code in the current directory:

```bash
holos generate platform k3d
```

Make the first commit:

```bash
git add .
git commit -m "holos generate platform k3d - $(holos --version)"
```

## Push the Platform Form

Push the Platform Form to the web app so we can provide top level configuration
values the platform components derive their final configuration from.

```bash
holos push platform form .
```

Visit the printed URL to view the Platform Form.

:::tip

You have complete control over the form fields and validation rules.

:::

## Submit the Platform Model

Fill out the form and submit the Platform Model.

:::tip

The default values are sufficient, however please create a public GitHub
repository to exercise ArgoCD functionality and update the Platform Model with
your `https://github.com` URL.

:::

## Pull the Platform Model

The Platform Model is the JSON representation of the Platform Form values.
Holos provides the Platform Model to CUE to render the platform configuration to
plain YAML. Configuration that varies is derived from the Platform Model using
CUE.

Pull the Platform Model to your local host to render the platform.

```bash
holos pull platform model .
```

The `platform.config.json` file is intended to be committed to version conrol.

```bash
git add platform.config.json
git commit -m "Add platform model"
```

:::warning

Do not store secret data using the Platform Form.  Holos uses ExternalSecret resources to securely sync with a Secret Store and ensure Secrets are never stored in version control.

:::

## Render the Platform

Rendering the platform iterates over each platform component and renders the
component into the final resources that will be sent to the API Server.

```bash
holos render platform ./platform
```

This command writes the fully rendered yaml to the `deploy/` directory.

:::warning

Do not edit the files in the `deploy` as they will be written over.

:::

You should have a tree similar to the following structure:

```txt
deploy
└── clusters
    └── workload
        ├── components
        │   ├── cert-manager
        │   │   └── cert-manager.gen.yaml
        │   └── namespaces
        │       └── namespaces.gen.yaml
        └── gitops
            ├── cert-manager.application.gen.yaml
            └── namespaces.application.gen.yaml
```

Commit the rendered platform configuration for `git diff` later.

```bash
git add deploy
git commit -m "holos render platform ./platform"
```

## Review the Platform Config

:::tip

This section is optional, included to provide insight into how Holos uses CUE
and Helm to unify and render the platform configuration.

:::

Take a moment to review the platform config `holos` rendered.

### ArgoCD Application

Note the Git URL you entered into the Platform Form is used to derive the ArgoCD
`Application` resource from the Platform Model.

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
    repoURL: https://github.com/holos-run/holos-k3d
    # highlight-next-line
    targetRevision: HEAD
```

One ArgoCD `Application` resource is produced for each Holos component by default.  Note the `cert-manger` component renders the output using Helm.  Holos unifies the Application resource using CUE.  The CUE definition which produces the rendered output is defined in `buildplan.cue` around line 222.

:::tip

Note how CUE does not use error-prone text templates, the language is well specified and typed which reduces errors when unifying the configuration with the Platform Model in the following `#Argo` definition.

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

Holos uses CUE to safely integrate the unmodified upstream `cert-manager` Helm
chart.

:::tip

Holos fully supports your existing Helm charts.  Consider leveraging `holos` as
an safer alternative to umbrella charts.

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

## Create the Workload Cluster

The Workload Cluster is where your applications and services will be deployed.
In production this is usually an EKS, GKE, or AKS cluster.

:::tip

Holos supports any compliant Kubernetes cluster and was developed and tested on GKE, EKS, Talos,
and Kubeadm clusters.

:::

```bash
k3d cluster create \
  --port "443:443@loadbalancer" \
  --k3s-arg "--disable=traefik@server:0" \
  workload
```

Traefik is disabled because Istio provides the same functionality.

## Apply the Platform Config

Use `kubectl` to apply each component to the cluster.  In production, it's common to fully automate this process with ArgoCD, but we use `kubectl` in development and exploration contexts to the same effect.

### Namespaces

```bash
kubectl apply --server-side=true -f ./deploy/clusters/workload/components/namespaces/namespaces.gen.yaml
```

### Cert Manager

```bash
kubectl apply --server-side=true -f ./deploy/clusters/workload/components/cert-manager/cert-manager.gen.yaml
```

The cert manager controller should start successfully:

```txt
❯ k get pods -A
NAMESPACE      NAME                                      READY   STATUS    RESTARTS   AGE
cert-manager   cert-manager-666cb4fb5f-dcp6x             1/1     Running   0          43s
cert-manager   cert-manager-cainjector-fd5479b67-jwbdr   1/1     Running   0          43s
cert-manager   cert-manager-webhook-588b7d86c8-rbcsm     1/1     Running   0          43s
kube-system    coredns-6799fbcd5-ksc2k                   1/1     Running   0          94m
kube-system    local-path-provisioner-6f5d79df6-b5js7    1/1     Running   0          94m
kube-system    metrics-server-54fd9b65b-mmx2n            1/1     Running   0          94m
```

### Gateway API

We use `HTTPRoute` resources from the Gateway API to expose services in a standard way.

```bash
kubectl apply --server-side=true -f ./deploy/clusters/workload/components/gateway-api/gateway-api.gen.yaml
```

### Istio Base

```bash
kubectl apply --server-side=true -f ./deploy/clusters/workload/components/istio-base/istio-base.gen.yaml
```

### Istio CNI

```bash
kubectl apply --server-side=true -f ./deploy/clusters/workload/components/istio-cni/istio-cni.gen.yaml
```

### Istio Controller

```bash
kubectl apply --server-side=true -f ./deploy/clusters/workload/components/istiod/istiod.gen.yaml
```

Once the Istio components have been applied, two pods should be running in the `istio-system` namespace:

```bash
kubectl get pods -n istio-system
```

```txt
NAME                      READY   STATUS    RESTARTS   AGE
istio-cni-node-hnvjg      0/1     Running   0          3m23s
istiod-6b48fd8448-lgmxd   1/1     Running   0          2m32s
```

### Istio Gateway

The Gateway component configures the ingress gateway used to expose services running in the cluster.

```bash
kubectl apply --server-side=true -f ./deploy/clusters/workload/components/gateway/gateway.gen.yaml
```

Once applied, the default Gateway Deployemnt should be running in the `istio-gateways` namespace.

```bash
kubectl get pods -n istio-gateways
```

```txt
NAME                             READY   STATUS    RESTARTS   AGE
default-istio-54f7fbd4f8-cl4v6   1/1     Running   0          8m3s
```


## Local CA

Create and install a local CA, necessary for secure https in your browser.

```bash
bash ./scripts/local-ca
```

:::note

Admin access is necessary to install the local CA into the system trust store.

:::

# Try Holos with k3d

Learn how to configure and deploy the Holos reference platform to your local host with k3d.

---

This guide assumes you are running the commands from your local host.
Capitalized terms have specific definitions described in the
[Glossary](/docs/glossary)

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
 2. [Docker](https://docs.docker.com/get-docker/) to use `k3d`.
 2. [holos](/docs/tutorial/install) to build the platform configuration.
 3. [kubectl](https://kubernetes.io/docs/tasks/tools/) to interact with the Kubernetes cluster.
 4. [helm](https://helm.sh/docs/intro/install/) to render Holos components that integrate vendor provided Helm charts.

## Install k3d

Refer to [k3d installation](https://k3d.io/#installation) to install `k3d`.

## Create the Workload Cluster

The Workload Cluster is where your applications and services will be deployed.
In production this is usually an EKS, GKE, or AKS cluster.  Holos supports any
compliant Kubernetes cluster and was developed and tested on GKE, EKS, Talos,
and Kubeadm clusters.

```bash
k3d cluster create --k3s-arg "--disable=traefik@server:0" workload
```

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

## Provide the Platform Model

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

Take a moment to review the platform config `holos` rendered.  Note the Git URL
you entered into the Platform Form is used to derive the ArgoCD `Application`
resource from the Platform Model.

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

## Apply the Platform Config

TODO

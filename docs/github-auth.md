# Kargo Demo

Kargo requires git credentials to promote artifacts.  Follow these steps to
setup you [Local Cluster] with these credentials.

## Process

First, fork this repo to your account.

We'll create a GitHub App, install the app with write permission to our own
fork, then store the private key in "$(mkcert -CAROOT)/kargo.yaml" so it's
automatically restored by the [reset-cluster] script.

### GitHub App

#### GitHub App Authentication

[Create a GitHub App](https://github.com/settings/apps/new) in the user or
organization where your bank-of-holos fork resides.

In the `GitHub App name` field, specify a unique name, for example `Holos - Local Cluster 1733418802` produced by:

```bash
echo -n "Holos - Local Cluster $(date +%s)" | pbcopy
```

Set the `Homepage URL` to `https://holos.run/docs/local-cluster/`.

Under `Webhook`, de-select `Active`.

Under `Permissions` → `Repository permissions` → `Contents`, select `Read and
write` permissions.  _The App will receive these permissions on all repositories
into which it is installed._

The `git-open-pr` step requires write permission to pull requests.  Add this
permission if you get the following error:

```
step execution failed: step 4 met error threshold of 1: failed to run step
"git-open-pr": error creating pull request: POST
https://api.github.com/repos/jeffmccune/kargo-demo/pulls: 403 Resource not
accessible by integration []
```

Under `Where can this GitHub App be installed?`, leave `Only on this account`
selected.

Click `Create GitHub App`.

Take note of the `App ID`. In your shell store it for use later using:

```bash
export BANK_OF_HOLOS_APP_ID=9999999
```

Scroll to the bottom of the page and click `Generate a private key`. The
resulting key will be downloaded immediately.  Record the path to this file for
use later using:

```bash
export BANK_OF_HOLOS_APP_KEY="$(ls -lr1 ~/Downloads/holos-local-cluster*.private-key.pem | tail -1)"
```

On the left-hand side of the page, click `Install App`.

Choose an account to install the App into by clicking `Install`.

Select `Only select repositories` and choose your `bank-of-holos` fork.
Remember that the App will receive the permissions you selected earlier for all
repositories you grant access.

Click `Install`.

In your browser's address bar, take note of the numeric identifier at the end of
the current page's URL. This is the `Installation ID`.  Save the installation id
for later.

For example, `https://github.com/settings/installations/99999999` is saved as:

```shell
export BANK_OF_HOLOS_INSTALL_ID=99999999
```

#### GitHub App Secret

Generate a Kubernetes Secret to store the Kargo git credentials.  We put this in
`mkcert -CAROOT` so `reset-cluster` restores it each time the local cluster is
reset.

Record the Git URL, the same as you set for `Organization.RepoURL`

```shell
export BANK_OF_HOLOS_REPO_URL="https://github.com/${USER}/bank-of-holos.git"
```

At this point you should have the following values, for example:

```shell
env | grep BANK_OF_HOLOS
```

```shell
BANK_OF_HOLOS_APP_ID=1079195
BANK_OF_HOLOS_APP_KEY=/Users/jeff/Downloads/holos-local-cluster-1733419264.2024-12-05.private-key.pem
BANK_OF_HOLOS_INSTALL_ID=58021430
BANK_OF_HOLOS_REPO_URL=https://github.com/jeffmccune/bank-of-holos.git
```

Generate the secret:

```shell
./scripts/kargo-git-creds
```

```txt
Secret created, apply with:
  kubectl apply -f '/Users/jeff/Library/Application Support/mkcert/kargo.yaml'

The reset-cluster script will automatically apply this secret going forward.
```

And apply it or reset your cluster.

```shell
kubectl apply -f '/Users/jeff/Library/Application Support/mkcert/kargo.yaml'
```

## Verification

Make sure you've configured Holos to use your `bank-of-holos` fork.

```shell
cat <<EOF > organization-repo-${USER}.cue
```
```cue showLineNumbers
@if($USER)
package holos

Organization: RepoURL: "${BANK_OF_HOLOS_REPO_URL}"
```
```shell
EOF
```

Then reset the cluster fully.  (Note this will delete and re-create your local
k3d cluster)

```bash
./scripts/full-reset
```

After a couple of minutes you should be able to log into https://kargo.holos.localhost with the admin password obtained with:

```shell
kubectl get secret -n kargo admin-credentials -o json \
  | jq --exit-status -r '.data.password | @base64d'
```

Make sure to commit to `main` and push it to your fork, then try and promote the
bank frontend.

ArgoCD is available at https://argocd.holos.localhost Most apps except those
which have previously been promoted in your fork should be in sync after a full
reset.

[Local Cluster]: https://holos.run/docs/local-cluster/
[reset-cluster]: ../scripts/reset-cluster

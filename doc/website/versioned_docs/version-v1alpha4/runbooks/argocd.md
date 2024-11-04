# ArgoCD

Create the deploy key secret in the management cluster.

```bash
tmp="$(mktemp -d)"
(cd $tmp && ssh-keygen -t ed25519 -f sshPrivateKey -m pem -C argocd -N '')
echo git@github.com:holos-run/holos-infra.git > "${tmp}/url"
holos create secret -n argocd --append-hash=false creds-holos-infra --from-file $tmp
rm -rf "$tmp"
```

When syncing the secret, the ExternalSecret needs to set the label
`argocd.argoproj.io/secret-type: repo-creds`.

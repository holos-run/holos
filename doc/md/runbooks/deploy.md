# Deployment

This document describes how deployment from `main` is configured.

1. Refer to the publish workflow.
2. Uses a SSH deploy key to:
3. Clone the holos-infra repo.
4. Write the image tag to saas/userdata/components/dev-holos-app/images.json
5. Run holos render platform ./platform
6. Commit and push the results.
7. ArgoCD takes over the rollout.

## Credentials

TODO: Lock this down more, the deploy key has too much access to the infra
repository.

```bash
mkdir -p tmp
cd tmp
ssh-keygen -t ed25519 -f holos-infra.key -m pem -C holos-run/holos -N ''
gh secret set DEPLOY_SSH_PRIVATE_KEY < holos-infra.key
gh api --method POST \
  -H "Accept: application/vnd.github+json" \
  /repos/holos-run/holos-infra/keys \
  -f title='holos-run/holos deploy key' \
  -f key="$(cat holos-infra.key.pub)" \
  -F read_only=false
cd ..
rm -rf tmp
```

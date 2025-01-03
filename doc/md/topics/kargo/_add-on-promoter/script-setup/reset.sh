set -xeuo pipefail
# DANGER MODE, don't reset the holos repo remotes...
cd kargo-demo
git remote add upstream https://github.com/holos-run/kargo-demo.git
git fetch upstream
git reset --hard upstream/main
git remote set-url origin git@github.com:${GH_USER}/kargo-demo.git
git push origin +HEAD:main

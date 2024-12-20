# Addon Promoter

This is a simple Kargo promotion pipeline for a cluster add-on.  It watches for
new helm chart versions, bumps the version number in the corresponding yaml
file, then submits a PR.

Tested with cert-manager and istio with Holos.

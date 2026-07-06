# Use Case: The Canonical Management Cluster

> Status: stub — owned by HOL-1506.

This chapter grounds the design in one canonical, end-to-end use case: a
management-cluster Profile composed of Roles (GitOps, secrets, ingress,
observability) whose Components are published CUE modules.  It walks the
full path — module `#Config` values set in a role, roles selected by the
management-cluster profile, the profile instantiated by a Platform, the
composed platform-wide TaskSet rendered to manifests — demonstrating that
Profile and Role remain pure CUE conventions while the Go tooling sees only
Platform → Component, and serving as the reference consumer for the
promotion gate described in [modules.md](modules.md).

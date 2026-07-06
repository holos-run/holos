# Rendering: Compiler Pool and the Platform-Wide DAG

> Status: stub — owned by HOL-1503.

This chapter specifies how v1beta1 renders a platform: every component's
TaskSet composes into one platform-wide TaskSet, and holos executes the
composed graph by topological sort over the DAG with a high level of
concurrency.  It covers the CUE compiler pool that amortizes evaluation cost
across components, task scheduling and failure semantics, cache keys for
hermetic task outputs, and how the platform-wide DAG subsumes the
per-component render loop of `internal/platform/` and `internal/component/`
without breaking v1alpha5/v1alpha6 components during migration.

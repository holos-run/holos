# Schema: TaskSet, Task, and Command

> Status: stub — owned by HOL-1502.

This chapter specifies the v1beta1 core schema: the `TaskSet` resource that
replaces the deprecated BuildPlan, the `Task` as the single unit of work in a
data transformation pipeline (subsuming the v1alpha5 Generator, Transformer,
and Validator concepts), and the first-class `Command` task for invoking
external tools.  It defines the struct-keyed (not list-keyed) field layout,
k8s-style lowerSnakeCase field naming, task input/output contracts, and the
dependency edges that make a TaskSet a DAG — the properties
[rendering.md](rendering.md) relies on to compose every component's TaskSet
into one platform-wide graph.

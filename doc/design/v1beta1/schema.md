# Schema: TaskSet, Task, and Command

> Status: stub — owned by HOL-1502.

This chapter specifies the v1beta1 core schema: the `TaskSet` resource that
replaces the deprecated BuildPlan, the `Task` as the single unit of work in a
data transformation pipeline (subsuming the v1alpha5 Generator, Transformer,
and Validator concepts), and the first-class `Command` task for invoking
external tools.  It defines the struct-keyed (not list-keyed) field layout,
the normative field-naming convention (resolving the v1alpha6 design note's
imprecise "k8s style lowerSnakeCase" wording in favor of the Kubernetes
lowerCamelCase API convention unless this chapter documents an intentional
departure), task input/output contracts, and the
dependency edges that make a TaskSet a DAG — the properties
[rendering.md](rendering.md) relies on to compose every component's TaskSet
into one platform-wide graph.

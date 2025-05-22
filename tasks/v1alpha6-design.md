# v1alpha6 Design

Design the new types for v1alpha6 learning from mistakes made in v1alpha5.

Mistakes:
- Lists are difficult to compose.  Use structs instead.
- Generators, Transformers, Validators are concepts.  Technically they're all Tasks in a data transformation pipeline.

Plan:  From the perspective of an user, start with the Platform entrypoint and work down to components.  The main design changes are:
1. Standardize on k8s style lowerSnakeCase for field names.
2. Replace lists with structs, e.g. Platform.spec.componets.
3. Deprecate BuildPlan, use a TaskSet instead.
4. Ensure TaskSets are composable into one big TaskSet for all platform components.
5. Execute the tasks in the big TaskSet using topological sort over the DAG with a high level of concurrency.

## TODO:

- [ ] Define the v1alpha6 Platform schema.
- [ ] Define the v1alpha6 TaskSet schema.
- [ ] Migrate the author schemas to use a core.TaskSet
- [ ] Helm
- [ ] Kustomize
- [ ] Resources
- [ ] File
- [ ] Command


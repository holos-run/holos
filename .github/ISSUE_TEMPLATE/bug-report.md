---
name: Bug report
about: Create a report to help us improve
title: ''
labels: NeedsInvestigation, Triage
assignees: ''
---

<!--
Please answer these questions before submitting your issue. Thanks!
To ask questions, see https://github.com/holos-run/holos/discussions
-->

### What version of holos are you using (`holos --version`)?

```
0.0.0
```

### Does this issue reproduce with the latest release?

<!--
Get the latest release with:

    brew install holos-run/tap/holos

Or see https://holos.run/docs/v1alpha5/tutorial/setup/
-->

### What did you do?

<!--
Please provide a testscript that should pass, but does not because of the bug.
See the below example.

You can create a txtar from a directory with:

  holos txtar ./path/to/dir

Refer to: https://github.com/rogpeppe/go-internal/tree/master/cmd/testscript
-->

Steps to reproduce:

```shell
testscript -v -continue <<EOF
```

```txtar
# Have: an error related to the imported Kustomize schemas.
# Want: holos show buildplans to work.
exec holos --version
exec holos init platform v1alpha5 --force
# remove the fix to trigger the bug
rm cue.mod/pkg/sigs.k8s.io/kustomize/api/types/var.cue
# want a BuildPlan shown
exec holos show buildplans
cmp stdout buildplan.yaml
# want this error to go away
! stderr 'cannot convert non-concrete value string'
-- buildplan.yaml --
kind: BuildPlan
-- platform/example.cue --
package holos

Platform: Components: example: {
	name: "example"
	path: "components/example"
}
-- components/example/example.cue --
package holos

import "encoding/yaml"

holos: Component.BuildPlan

Component: #Kustomize & {
	KustomizeConfig: Kustomization: patches: [
		{
			target: kind: "CustomResourceDefinition"
			patch: yaml.Marshal([{
				op:    "add"
				path:  "/metadata/annotations/example"
				value: "example-value"
			}])
		},
	]
}
```
```shell
EOF
```

### What did you expect to see?

The testscript should pass.

### What did you see instead?

The testscript fails because of the bug.

```txt
# Have: an error related to the imported Kustomize schemas.
# Want: holos show buildplans to work. (0.168s)
> exec holos --version
[stdout]
0.100.1-2-g9b10e23-dirty
> exec holos init platform v1alpha5 --force
# remove the fix to trigger the bug (0.000s)
> rm cue.mod/pkg/sigs.k8s.io/kustomize/api/types/var.cue
# want a BuildPlan shown (0.091s)
> exec holos show buildplans
[stderr]
could not run: holos.spec.artifacts.0.transformers.0.kustomize.kustomization.patches.0.target.name: cannot convert non-concrete value string at builder/v1alpha5/builder.go:218
holos.spec.artifacts.0.transformers.0.kustomize.kustomization.patches.0.target.name: cannot convert non-concrete value string:
    $WORK/cue.mod/gen/sigs.k8s.io/kustomize/api/types/var_go_gen.cue:33:2
[exit status 1]
FAIL: <stdin>:8: unexpected command failure
> cmp stdout buildplan.yaml
diff stdout buildplan.yaml
--- stdout
+++ buildplan.yaml
@@ -0,0 +1,1 @@
+kind: BuildPlan

FAIL: <stdin>:9: stdout and buildplan.yaml differ
# want this error to go away (0.000s)
> ! stderr 'cannot convert non-concrete value string'
FAIL: <stdin>:11: unexpected match for `cannot convert non-concrete value string` found in stderr: cannot convert non-concrete value string
failed run
```

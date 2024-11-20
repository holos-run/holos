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

```shell
holos --version
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
brew install testscript
```

```shell
testscript -v -continue example.txt
```

```txt
# Have: an error related to the imported Kustomize schemas.
# Want: holos show buildplans to work.
exec holos --version
exec holos init platform v1alpha5 --force
# want a BuildPlan shown
exec holos show buildplans
stdout 'kind: BuildPlan'
# want this error to go away
! stderr 'cannot convert non-concrete value string'
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

### What did you expect to see?

The testscript should pass.

### What did you see instead?

The testscript fails because of the bug.

```shell
testscript -v -continue example.txt
```

```txt
# Have: an error related to the imported Kustomize schemas.
# Want: holos show buildplans to work. (0.073s)
> exec holos --version
[stdout]
0.100.0
> exec holos init platform v1alpha5 --force
# want a BuildPlan shown (0.085s)
> exec holos show buildplans
[stderr]
could not run: holos.spec.artifacts.0.transformers.0.kustomize.kustomization.patches.0.target.name: cannot convert non-concrete value string at builder/v1alpha5/builder.go:218
holos.spec.artifacts.0.transformers.0.kustomize.kustomization.patches.0.target.name: cannot convert non-concrete value string:
    $WORK/cue.mod/gen/sigs.k8s.io/kustomize/api/types/var_go_gen.cue:33:2
[exit status 1]
FAIL: example.txt:7: unexpected command failure
> stdout 'kind: BuildPlan'
FAIL: example.txt:8: no match for `kind: BuildPlan` found in stdout
# want this error to go away (0.000s)
> ! stderr 'cannot convert non-concrete value string'
FAIL: example.txt:10: unexpected match for `cannot convert non-concrete value string` found in stderr: cannot convert non-concrete value string
failed run
```

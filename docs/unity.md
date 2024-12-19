# Unity

Test cases in this repository should work with Unity now.  Dependencies on
external executables have been eliminated provided we test against `holos show
platform` and `holos show buildplans`.  These two commands exercise CUE without
executing the resulting build plans.

Verified with the following command at the root:

```bash
docker run -v $(pwd):/src -w /src golang:1.23 go test
```

```
PASS
ok      github.com/holos-run/kargo-demo 0.040s
```

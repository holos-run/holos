# Unity

## Integration Test

To start, please execute the following to see if this demo repo produces all of
the BuildPlan resources we need:

```bash
make install
make unity
```

This is equivalent to:

```bash
go install github.com/holos-run/holos/cmd/holos

# this should work
export CUE_EXPERIMENT=evalv3=0
holos show buildplans

# this should also work but probably does not
export CUE_EXPERIMENT=evalv3=1
holos show buildplans
```

## Test Scripts

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

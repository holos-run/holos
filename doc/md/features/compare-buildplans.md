# Compare Buildplans

Use the `holos compare buildplans <f1> <f2>` command to compare two BuildPlan
Files.  Useful to ensure different configuration versions produce the same
results.

The `holos show buildplans` command writes a BuildPlan File to standard output.
A BuildPlan File is a yaml encoded stream of BuildPlan objects.

## User Requirements
 1. `holos compare buildplans before.yaml after.yaml` must return exit code 1 when after.yaml contains fields (recursively) not present in before.yaml
 2. `holos compare buildplans before.yaml after.yaml --backwards-compatible` must return exit code 0 when after.yaml contains fields (recursively) not present in before.yaml

## Behavior Specification
BuildPlan File f1 is equivalent to f2 when:
 1. f1 and f2 have an equal number of BuildPlan objects.
 2. each object in f1 is equivalent to exactly one unique object in f2.

Two BuildPlans, bp1 and bp2, are equivalent when:
 1. All field values in bp1 are equivalent to the same field in bp2.
 2. Both 1 and 2 apply to nested objects, recursively.
 3. Field f is equivalent when bp1.f exactly equals bp2.f, except for:
    3.1. Objects in the spec.artifacts list may appear in any arbitrary order.
    3.2. The ordering of keys does not matter.
 4. Backwards compatibility behavior (controlled by isBackwardsCompatible):
    - When false: bp2 and bp1 must have exactly the same fields
    - When true: bp2 may have additional fields that don't exist in bp1
    (e.g., new features added in a newer version)
    Example:
    bp1 has {name: "x", version: "1.0"}
    bp2 has  {name: "x", version: "1.0", newFeature: "enabled"}
    This comparison passes when isBackwardsCompatible=true
 5. Fields in bp1 must always be present in bp2 (regardless of backwards
    compatibility mode).
 6. List type fields with a null value are equivalent to:
    6.1. null values
    6.2. empty values ([])
    6.2. a missing field

A BuildPlan File is valid when:
 1. Two or more identical objects exist in the same file.  They must be
    treated as unique objects when comparing BuildPlan Files
 2. Two objects may have the same value for the metadata.name field.
 3. The kind field of all objects in the file stream is "BuildPlan"

## Implementation Guidance

1. Implement a stub Comparer struct in the internal/compare package.
2. Implement a stub Comparer.BuildPlans() method.
3. Write test cases for each item in the Behavior Specification section.  Use a table test approach that loads each test case from a subdirectory and reads the test case data from a `testcase.json` file.  The json file should have an exitCode, name, msg, file1 and file2 fields.  file1 is "before.yaml" and file2 is "after.yaml".
4. Modify the Comparer.BuildPlans() method to satisfy each test case.
5. Using the existing commands as an example, wire up the command line to the compare package.

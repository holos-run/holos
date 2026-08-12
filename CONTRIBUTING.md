# Contributing to Holos

Thank you for your interest in contributing to Holos! We welcome all
contributions — code, documentation, bug reports, feature proposals, and
experience reports from using Holos to manage your platform. This guide
explains how to get involved and what to expect along the way.

## Ways to Contribute

You do not need to write code to make a meaningful contribution. You can:

- Report bugs or propose enhancements in the
  [issue tracker](https://github.com/holos-run/holos/issues).
- Ask and answer questions in
  [GitHub Discussions](https://github.com/holos-run/holos/discussions) or the
  [Discord server](https://discord.gg/JgDVbNpye7).
- Improve the [documentation](https://holos.run/docs/overview/), tutorials,
  and examples.
- Share an experience report describing how you use Holos and where it falls
  short — these directly shape the roadmap.
- Contribute code: bug fixes, features, and tests.

## Before You Start

### Sign the Contributor License Agreement

Before we can accept your contributions, you must agree to the
[Contributor License Agreement](CLA.md) (CLA). The CLA clarifies the
intellectual property license you grant with your contributions. It protects
you as a contributor as well as Open Infrastructure Operations LLC and users
of the software; it does not change your rights to use your own contributions
for any other purpose.

If your employer has rights to intellectual property you create, make sure
you have permission to contribute on their behalf, or that your employer has
waived such rights for your contributions.

### Discuss Your Change First

For anything beyond a trivial fix (typos, small documentation corrections),
please open an issue before writing code. The issue tracker is the main form
of communication between contributors, and discussing the design up front
avoids wasted effort — code review is for reviewing implementation, not for
debating whether a change should happen at all.

Significant changes to the API surface (`/api/core/` and `/api/author/`)
warrant a design discussion in the issue before implementation begins.

## Development Environment

You need [Go](https://go.dev/dl/) (see `go.mod` for the minimum version) and
`make`. Clone the repository and verify your environment:

```bash
git clone https://github.com/holos-run/holos.git
cd holos
make test
```

Common development targets:

```bash
make build         # Build the holos binary
make install       # Install holos to GOPATH/bin (required to test CLI changes)
make test          # Run all tests
make lint          # Run go vet and golangci-lint
make fmt           # Format Go code
make coverage      # Generate a test coverage report
```

Run `make help` to see the full list of targets.

## Making Changes

### Branch and Commit

Create a topic branch from `main` for your change. Keep the commit history of
your pull request minimal — squash fixup commits before marking the PR ready
for review, and rebase on `main` rather than merging it in.

### Commit Messages

Follow the Go project's commit message convention. The first line is a short
summary prefixed with the name of the primary package or area affected, in
lower case, with no trailing period:

```
component: retry chart pulls on transient network errors

Helm chart pulls previously failed permanently on the first network
error. Retry with exponential backoff up to three attempts so transient
failures during render do not fail the whole platform.

Fixes #123
```

Conventions:

- The first line completes the sentence "This change modifies Holos to ___".
- Separate the summary from the body with a blank line.
- Write the body in complete sentences, wrapped at about 72 columns, and
  explain *why* the change is needed, not just what it does.
- Reference related issues with `Fixes #123` (closes the issue on merge) or
  `Updates #123` (leaves it open).

### Code Style

- Format Go code with `make fmt` and CUE files with `cue fmt`.
- Prefer `errors.Format()` from `/internal/errors/` over `fmt.Errorf()`.
- Use structured logging with `slog`; obtain the logger with
  `logger.FromContext(ctx)`.
- Follow the existing Cobra command patterns in `/internal/cli/` for CLI
  changes.
- Develop new functionality against the in-development API version (currently
  v1alpha6); v1alpha5 is stable and changes to it should be limited to bug
  fixes.
- Avoid adding external dependencies unless they are strictly necessary.

### Testing

All changes must pass `make test` and `make lint`. New features will not be
accepted without suitable test coverage, and bug fixes should include a
regression test that fails without the fix.

- Unit tests live in `*_test.go` files colocated with the source.
- Integration tests live in `/cmd/holos/tests/`.
- Example platforms used as test fixtures live in
  `/internal/testutil/fixtures/`.
- Run a single test with `go test -run TestName ./path/to/package`.

## Submitting a Pull Request

1. Confirm you have agreed to the [CLA](CLA.md).
2. Rebase your branch on the current `main`.
3. Run `make lint` and `make test` locally.
4. Open the pull request as a **draft** while work is in progress; mark it
   ready for review when it is complete. This reduces notification noise for
   maintainers.
5. Give the PR a clear title following the commit message convention, and
   describe what the change does, why it is needed, and link the related
   issue.

## The Review Process

Expect review comments — they are a normal and healthy part of the process,
and multiple review rounds are common even for experienced contributors.
Treat each review comment like a ticket you are expected to close: either
make the requested change and reply "Done", or explain why you chose a
different approach. Don't be discouraged by iteration; it is how the code
gets better.

Once a maintainer approves your PR and CI passes, a maintainer will merge it.

## Getting Help

If you get stuck at any point:

- Ask in the [Discord server](https://discord.gg/JgDVbNpye7).
- Start a [GitHub Discussion](https://github.com/holos-run/holos/discussions).
- Comment on the relevant issue.

Thank you for contributing!

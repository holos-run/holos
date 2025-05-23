# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Holos is a configuration management tool for Kubernetes that implements the rendered manifests pattern using CUE. It unifies Helm charts, Kustomize bases, and raw Kubernetes manifests into a single, declarative pipeline.

### Core Flow
```
Platform → Components → BuildPlan → Generators → Transformers → Validators → Manifests
```

## Key Commands

```bash
# Development
make build         # Build the binary
make install       # Install binary (REQUIRED before testing holos commands)
make test          # Run all tests
make fmt           # Format Go code
make lint          # Run linters
make coverage      # Generate coverage report

# Documentation
make update-docs   # Update generated docs
make website       # Build the documentation website

# Usage (run 'make install' first to test code changes)
holos render platform    # Render entire platform
holos render component   # Render single component
holos show buildplans    # Show build plans
holos init platform      # Initialize new platform
```

## Architecture

### Directory Structure
- `/api/` - API definitions (v1alpha5 stable, v1alpha6 in development)
- `/cmd/` - CLI entry point
- `/internal/cli/` - Command implementations
- `/internal/component/` - Component handling logic
- `/internal/platform/` - Platform handling logic
- `/internal/generate/` - Code generation

### Key Files
- `/internal/cli/render/render.go` - Core render logic
- `/internal/component/component.go` - Component processing
- `/api/core/v1alpha*/types.go` - API type definitions

### Component Types
1. **Helm** - Wraps Helm charts
2. **Kustomize** - Wraps Kustomize bases  
3. **Kubernetes** - Raw Kubernetes manifests

## CUE Patterns

Components are defined in CUE:
```cue
package holos

holos: Component.BuildPlan

Component: #Helm & {
    Name: "example"
    Chart: {
        version: "1.0.0"
        repository: {
            name: "example"
            url:  "https://charts.example.com"
        }
    }
}
```

## Testing

- Unit tests: `*_test.go` files colocated with source
- Integration tests: `/cmd/holos/tests/`
- Example platforms: `/internal/testutil/fixtures/`
- Run single test: `go test -run TestName ./path/to/package`

## Development Patterns

1. Error handling: Prefer `errors.Format()` from `/internal/errors/` over `fmt.Errorf()`
2. Logging: Use structured `slog`, get logger with `logger.FromContext(ctx)`
3. CLI commands: Follow Cobra patterns in `/internal/cli/`
4. CUE formatting: Always run `cue fmt` on CUE files
5. Go formatting: Always run `go fmt` on go files
6. Develop against v1alpha6 packages.
7. Commits: Use the package name as the first word in the commit, lower case.  Commit without asking permission.  Always run `make lint` and `make test` before committing.

## Version Management

- Version files: `/version/embedded/{major,minor,patch}`
- Bump version: `make bump`
- API versions: v1alpha5 (stable), v1alpha6 (development)

## Key Concepts

- **Platform**: Top-level configuration containing all components
- **Component**: Unit of configuration (DAG of Tasks producing deployment configs for one component)
- **TaskSet**: DAG of Tasks (Similar to how make tasks behave)
- **BuildPlan**: Instructions for building a component.  Deprecated in v1alpha6, use TaskSet instead.
- **Generator**: Creates manifests (Helm, Kustomize, etc.) author schema only in v1alpha6
- **Transformer**: Modifies generated manifests, author schema only in v1alpha6
- **Validator**: Validates final manifests, author schema only in v1alpha6

## Resources

- Tutorials: `/doc/md/tutorial/`
- Platform templates: `/internal/generate/platforms/`
- Test fixtures: `/internal/testutil/fixtures/`
- Core schemas: `/api/core/` (Abstraction over low level data pipeline tasks)
- Author schemas: `/api/author/` (User facing abstractions over core Schemas)
- Task planning documents are located in the `/tasks/` directory

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
make test          # Run all tests
make fmt           # Format Go code
make lint          # Run linters
make coverage      # Generate coverage report

# Documentation
make update-docs   # Update generated docs
make website       # Build the documentation website

# Usage
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

1. Error handling: Use `internal/errors/` types, wrap with context
2. Logging: Use structured `slog`, get logger with `logger.FromContext(ctx)`
3. CLI commands: Follow Cobra patterns in `/internal/cli/`
4. CUE formatting: Always run `cue fmt` on CUE files
5. Develop against v1alpha6 packages.

## Version Management

- Version files: `/version/embedded/{major,minor,patch}`
- Bump version: `make bump`
- API versions: v1alpha5 (stable), v1alpha6 (development)

## Key Concepts

- **Platform**: Top-level configuration containing all components
- **Component**: Unit of configuration (Helm/Kustomize/Kubernetes)
- **BuildPlan**: Instructions for building a component
- **Generator**: Creates manifests (Helm, Kustomize, etc.)
- **Transformer**: Modifies generated manifests
- **Validator**: Validates final manifests

## Resources

- Tutorials: `/doc/md/tutorial/`
- Platform templates: `/internal/generate/platforms/`
- Test fixtures: `/internal/testutil/fixtures/`

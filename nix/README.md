# Nix Development Environment

This directory contains configuration pertaining to building and developing the
Holos CLI with Nix. The main [Nix flake](https://wiki.nixos.org/wiki/Flakes)
configuration is defined in the repository root at `./flake.nix`, and all
commands in this README assume they'll be run from the repository root.

## Bootstrapping

If you don't have Nix installed already, you can review commands required to do
so along with direnv (to auto-activate the development environment) by running:

```bash
make -n -f nix/nix.mk bootstrap
```

This would:

1. Install Nix using the Determinate Systems installer
2. Install direnv for automatic environment switching
3. Guide you through setting up your shell

Remove the `-n` if you understand the commands that will be executed.

## Common Commands

### Development

```bash
# Enter development shell
nix develop

# Format Nix files
nix fmt

# Show available development shell and packages
nix flake show

# Show flake inputs and their versions
nix flake metadata
```

### Building

```bash
# Build default package (holos CLI)
nix build

# Build specific output
nix build .#<output-name>
```

## Directory Structure

```bash
nix/
├── devshell.nix   # Development environment configuration
├── holos.nix      # Main package derivation
└── nix.mk         # Bootstrap makefile
```

## Environment Management

You can use either direnv or nix develop to manage your development environment:

```bash
# Using direnv (automatic)
cd /path/to/project    # Environment activates automatically
direnv allow          # First time only
direnv revoke         # Disable automatic activation

# Using nix develop (manual)
nix develop           # Enter development environment
exit                  # Leave development environment
```

## Updating Dependencies

To update Nix dependencies:

```bash
# Update all inputs
nix flake update

# Update specific input
nix flake lock --update-input nixpkgs
```

For more information about Nix, see:

- [Nix Manual](https://nixos.org/manual/nix/stable/)
- [Nix Flakes](https://nixos.wiki/wiki/Flakes)

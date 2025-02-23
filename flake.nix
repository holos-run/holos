{
  description = "Holos CLI tool";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    systems.url = "github:nix-systems/default";
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import inputs.systems;

      perSystem =
        {
          config,
          self',
          inputs',
          pkgs,
          system,
          ...
        }:
        let
          version = "0.104.1"; # Should match holos.nix
        in
        {
          devShells = {
            default = pkgs.callPackage ./nix/devshell.nix {
              inherit (self') packages;
            };
          };

          packages = {
            default = pkgs.callPackage ./nix/holos.nix { };

          };

          formatter = pkgs.nixfmt-rfc-style;
        };
    };
}

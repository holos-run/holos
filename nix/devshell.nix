{
  pkgs,
  packages,
}:
pkgs.mkShell {
  inputsFrom = [ packages.default ];
  # see https://search.nixos.org/ for package name ids
  packages = with pkgs; [
    cue
    go
    gopls
    k3d
    kind
    kubectl
    kubectx # includes kubens
    k9s
    kubernetes-helm
    kustomize
    timoni
    packages.default
  ];
}

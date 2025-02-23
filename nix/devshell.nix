{
  pkgs,
  packages,
}:
pkgs.mkShell {
  inputsFrom = [packages.default];
  packages = with pkgs; [
    go
    gopls
    packages.default
  ];
}

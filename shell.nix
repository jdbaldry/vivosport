{ pkgs ? import <nixpkgs> }:
with pkgs;
mkShell {
  buildInputs = [
    gnumake
  ] ++ [
    go-outline
    go-tools
    go_1_16
    goimports
    golangci-lint
    gopkgs
    gopls
    pgformatter
    sqlc
  ];
  shellHook = ''
    # ...
  '';
}

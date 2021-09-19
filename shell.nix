{ pkgs ? import <nixpkgs> }:
with pkgs;
mkShell {
  buildInputs = [
    docker-compose
    fitgen
    gnumake
    postgresql
    sqlc
  ] ++ [
    go-outline
    go-tools
    go_1_16
    goimports
    golangci-lint
    gopkgs
    gopls
  ];
  shellHook = ''
    export PATH="$PATH:$(pwd)/result/bin"
  '';
}

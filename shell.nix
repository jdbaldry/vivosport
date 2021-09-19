{ pkgs ? import <nixpkgs> }:
with pkgs;
mkShell {
  buildInputs = [
    curl
    docker-compose
    fitgen
    libreoffice
    gnumake
    go-jsonnet
    mount
    openjdk
    postgresql
    rsync
    sqlc
    unzip
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

{
  description = "vivosport shell development tooling";

  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    {
      overlay = import ./overlay.nix;
    } //
    (flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs
          {
            inherit system; overlays = [ self.overlay ];
          };
      in
      {
        defaultPackage = import ./default.nix { inherit pkgs; };
        devShell = import ./shell.nix { inherit pkgs; };
      }));
}

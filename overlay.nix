final: prev:

rec {
  sqlc = prev.callPackage ./sqlc.nix { pkgs = prev; };
}

final: prev:

rec {
  fitgen = prev.callPackage ./fitgen.nix { pkgs = prev; };
  sqlc = prev.callPackage ./sqlc.nix { pkgs = prev; };
}

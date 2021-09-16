{ pkgs ? import <nixpkgs> }:

with pkgs;
buildGoModule {
  pname = "vivosport";
  version = "0.0.1";

  meta = with lib; {
    maintainers = with maintainers; [ jdbaldry ];
  };
  src = lib.cleanSource ./.;
  vendorSha256 = null;
}

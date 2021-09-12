{ pkgs ? import <nixpkgs> }:

with pkgs;
buildGoModule {
  pname = "fit";
  version = "0.0.1";

  meta = with lib; {
    maintainers = with maintainers; [ jdbaldry ];
  };
  src = lib.cleanSource ./.;
  vendorSha256 = "sha256-pQpattmS9VmO3ZIQUFn66az8GSmB4IvYhTTCFn6SUmo=";
}

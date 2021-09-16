{ pkgs ? import <nixpkgs> }:

with pkgs;
buildGoModule rec {
  pname = "sqlc";
  version = "1.10.0";

  buildInputs = [ xxHash ];
  doCheck = false;
  meta = with lib; {
    maintainers = with maintainers; [ jdbaldry ];
  };
  src = fetchFromGitHub {
    owner = "kyleconroy";
    repo = pname;
    rev = "v${version}";
    sha256 = "sha256-tGkqckCCDQX4X4i/pgzLt+EKzytK1dsld85Odfs05To=";
  };
  runVend = true;
  vendorSha256 = "sha256-gxzmWjhGXACPLyOrquoCw6XN1vKqXDh7WrsYkxpHYkw=";
}

{ pkgs ? import <nixpkgs> }:

with pkgs;
buildGoModule rec {
  pname = "fit";
  version = "0.10.0";

  meta = with lib; {
    maintainers = with maintainers; [ jdbaldry ];
  };
  src = fetchFromGitHub {
    owner = "tormoder";
    repo = pname;
    rev = "v${version}";
    sha256 = "sha256-QBVfmedQ6bPeHosZAICPoqNJqSIpgtXSkwnh6pzYcL4=";
  };
  vendorSha256 = "sha256-sP4ZpGJ3q05D/LFIX4RjkBPtO/7QmWeitHl0irSny6A=";
}

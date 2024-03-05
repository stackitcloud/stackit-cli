#{ lib, buildGoModule, fetchFromGitHub }:
{ pkgs ? import <nixpkgs> {} }:
with pkgs;

let
  version = "0.1.0-beta.8";
  #date = toString builtins.currentTime;
in

buildGoModule {
  pname = "stackit-cli";
  inherit version;

  src = fetchFromGitHub {
    owner = "stackitcloud";
    repo = "stackit-cli";
    rev = "v${version}";
    hash = "sha256-YJLtGrgYXuExR8Q6f25HW5ovJKAm2rfZ90qCWQHrZLc=";
  };
  vendorHash = "sha256-u/rLY03utuD98IPKbdeFrBuCmhZ0vbJ+3tlmAo1V1GU=";

  CGO_ENABLED = 0;
  ldflags = [ "-X main.version=${version}" ];
  excludedPackages = [ "scripts" ];

  nativeCheckInputs = [ less ];
  preCheck = ''
    export HOME=$TMPDIR
  '';

  postInstall = ''
    mv $out/bin/{stackit-cli,stackit} # rename the binary
  '';

  # meta = with lib; {
  #   description = "STACKIT CLI";
  #   homepage = "https://github.com/stackitcloud/stackit-cli";
  #   license = licenses.asl20;
  #   # maintainers = with maintainers; [];
  # };
}
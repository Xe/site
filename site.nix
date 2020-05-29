{ pkgs ? import (import ./nix/sources.nix).nixpkgs }:
with pkgs;

assert lib.versionAtLeast go.version "1.13";

buildGoPackage rec {
  name = "christinewebsite-HEAD";
  version = "latest";
  goPackagePath = "christine.website";
  src = ./.;
  goDeps = ./nix/deps.nix;
  allowGoReference = false;

  preBuild = ''
    export CGO_ENABLED=0
    buildFlagsArray+=(-pkgdir "$TMPDIR")
  '';

  postInstall = ''
    cp -rf $src/blog $out/blog
    cp -rf $src/css $out/css
    cp -rf $src/gallery $out/gallery
    cp -rf $src/signalboost.dhall $out/signalboost.dhall
    cp -rf $src/static $out/static
    cp -rf $src/talks $out/talks
    cp -rf $src/templates $out/templates
  '';
}

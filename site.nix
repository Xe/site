{ sources ? import ./nix/sources.nix, pkgs ? import sources.nixpkgs { } }:
with pkgs;

let
  srcNoTarget = dir:
    builtins.filterSource
    (path: type: type != "directory" || builtins.baseNameOf path != "target")
    dir;
  naersk = pkgs.callPackage sources.naersk { };
  gruvbox-css = pkgs.callPackage sources.gruvbox-css { };
  src = srcNoTarget ./.;
  xesite = naersk.buildPackage {
    inherit src;
    buildInputs = [ pkg-config openssl ];
    remapPathPrefix = true;
  };

in pkgs.stdenv.mkDerivation {
  inherit (xesite) name;
  inherit src;
  phases = "installPhase";

  installPhase = ''
    mkdir -p $out $out/blog $out/css $out/gallery $out/static $out/talks $out/bin

    cp -rf $src/config.dhall $out/config.dhall
    cp -rf $src/blog $out/blog
    cp -rf $src/css $out/css
    cp -rf $src/gallery $out/gallery
    cp -rf $src/signalboost.dhall $out/signalboost.dhall
    cp -rf $src/static $out/static
    cp -rf $src/talks $out/talks

    cp -rf ${xesite}/bin/xesite $out/bin/xesite
  '';
}

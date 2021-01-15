{ sources ? import ./nix/sources.nix, pkgs ? import sources.nixpkgs { } }:
with pkgs;

let
  rust = pkgs.callPackage ./nix/rust.nix { };

  srcNoTarget = dir:
    builtins.filterSource
    (path: type: type != "directory" || builtins.baseNameOf path != "target")
    dir;

  naersk = pkgs.callPackage sources.naersk {
    rustc = rust;
    cargo = rust;
  };
  dhallpkgs = import sources.easy-dhall-nix { inherit pkgs; };
  src = srcNoTarget ./.;

  xesite = naersk.buildPackage {
    inherit src;
    doCheck = true;
    buildInputs = [ pkg-config openssl git ];
    remapPathPrefix = true;
  };

  config = stdenv.mkDerivation {
    pname = "xesite-config";
    version = "HEAD";
    buildInputs = [ dhallpkgs.dhall-simple ];

    phases = "installPhase";

    installPhase = ''
      cd ${src}
      dhall resolve < ${src}/config.dhall >> $out
    '';
  };

in pkgs.stdenv.mkDerivation {
  inherit (xesite) name;
  inherit src;
  phases = "installPhase";

  installPhase = ''
    mkdir -p $out $out/bin

    cp -rf ${config} $out/config.dhall
    cp -rf $src/blog $out/blog
    cp -rf $src/css $out/css
    cp -rf $src/gallery $out/gallery
    cp -rf $src/signalboost.dhall $out/signalboost.dhall
    cp -rf $src/static $out/static
    cp -rf $src/talks $out/talks

    cp -rf ${xesite}/bin/xesite $out/bin/xesite
  '';
}

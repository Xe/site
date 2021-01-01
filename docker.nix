{ system ? builtins.currentSystem }:

let
  sources = import ./nix/sources.nix;
  pkgs = import sources.nixpkgs { inherit system; };
  callPackage = pkgs.lib.callPackageWith pkgs;
  site = callPackage ./site.nix { };

  dockerImage = pkg:
    pkgs.dockerTools.buildLayeredImage {
      name = "xena/christinewebsite";
      tag = "latest";

      contents = [ pkgs.cacert pkg ];

      config = {
        Cmd = [ "${pkg}/bin/xesite" ];
        Env = [ "CONFIG_FNAME=${pkg}/config.dhall" "RUST_LOG=info" ];
        WorkingDir = "/";
      };
    };

in dockerImage site

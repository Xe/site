{ system ? builtins.currentSystem }:

let
  pkgs = import (import ./nix/sources.nix).nixpkgs { inherit system; };
  callPackage = pkgs.lib.callPackageWith pkgs;
  site = callPackage ./site.nix { };

  dockerImage = pkg:
    pkgs.dockerTools.buildLayeredImage {
      name = "xena/christinewebsite";
      tag = "latest";

      contents = [ pkgs.cacert ];

      config = {
        Cmd = [ "${pkg}/bin/xesite" ];
        WorkingDir = "${pkg}/";
      };
    };

in dockerImage site

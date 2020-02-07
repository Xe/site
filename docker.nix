{ system ? builtins.currentSystem }:

let
  pkgs = import <nixpkgs> { inherit system; };

  callPackage = pkgs.lib.callPackageWith pkgs;

  site = callPackage ./site.nix { };

  dockerImage = pkg:
    pkgs.dockerTools.buildImage {
      name = pkg.name;
      tag = pkg.version;

      contents = [ pkg ];

      config = {
        Cmd = [ "/bin/site" ];
        WorkingDir = "/";
      };
    };

in dockerImage site

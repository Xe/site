{ }:

let
  sources = import ./nix/sources.nix;
  pkgs = import sources.nixpkgs { };
in pkgs.callPackage ./site.nix { pkgs = pkgs; }

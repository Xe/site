let
  sources = import ./sources.nix;
  pkgs = import sources.nixpkgs { };
  dhall = import sources.easy-dhall-nix { inherit pkgs; };
in dhall.dhall-yaml-simple

{ pkgs ? import <nixpkgs> {} }:

pkgs.callPackage ./site.nix {}

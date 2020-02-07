let
  pkgs = import <nixpkgs> { };
  sources = import ./nix/sources.nix;
in pkgs.mkShell { buildInputs = [ pkgs.go sources.vgo2nix sources.niv ]; }

let
  pkgs = import <nixpkgs> { };
  sources = import ./nix/sources.nix;
  vgo2nix = (import sources.vgo2nix { });
in pkgs.mkShell { buildInputs = [ pkgs.go pkgs.niv vgo2nix ]; }

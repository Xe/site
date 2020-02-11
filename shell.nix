let
  pkgs = import <nixpkgs> { };
  sources = import ./nix/sources.nix;
  xepkgs = import sources.xepkgs { };
  vgo2nix = import sources.vgo2nix { };
in pkgs.mkShell { buildInputs = [ pkgs.go pkgs.niv xepkgs.gopls vgo2nix ]; }

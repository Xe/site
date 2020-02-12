let
  sources = import ./nix/sources.nix;
  pkgs = import sources.nixpkgs { };
  xepkgs = import sources.xepkgs { };
  vgo2nix = import sources.vgo2nix { };
in pkgs.mkShell { buildInputs = [ pkgs.go pkgs.niv xepkgs.gopls vgo2nix ]; }

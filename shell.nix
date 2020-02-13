let
  sources = import ./nix/sources.nix;
  niv = (import sources.niv { }).niv;
  pkgs = import sources.nixpkgs { };
  xepkgs = import sources.xepkgs { };
  vgo2nix = import sources.vgo2nix { };
in pkgs.mkShell { buildInputs = [ pkgs.go xepkgs.gopls niv vgo2nix ]; }

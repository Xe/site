let
  sources = import ./nix/sources.nix;
  pkgs = import sources.nixpkgs { };
  niv = (import sources.niv { }).niv;
  dhallpkgs = import sources.easy-dhall-nix { inherit pkgs; };
  dhall-yaml = dhallpkgs.dhall-yaml-simple;
  dhall = dhallpkgs.dhall-simple;
  xepkgs = import sources.xepkgs { inherit pkgs; };
  vgo2nix = import sources.vgo2nix { inherit pkgs; };
in with pkgs;
with xepkgs;
mkShell {
  buildInputs = [
    # Go tools
    go
    goimports
    gopls
    vgo2nix

    # kubernetes deployment
    dhall
    dhall-yaml

    # dependency manager
    niv
  ];
}

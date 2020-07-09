let
  sources = import ./nix/sources.nix;
  pkgs = import sources.nixpkgs { };
  dhallpkgs = import sources.easy-dhall-nix { inherit pkgs; };
  dhall-yaml = dhallpkgs.dhall-yaml-simple;
  dhall = dhallpkgs.dhall-simple;
  xepkgs = import sources.xepkgs { inherit pkgs; };
in with pkgs;
with xepkgs;
mkShell {
  buildInputs = [
    # Go tools
    go
    goimports
    gopls
    vgo2nix

    # Rust
    cargo
    cargo-watch
    rls
    rustc
    rustfmt

    # kubernetes deployment
    dhall
    dhall-yaml

    # dependency manager
    niv

    # tools
    ispell
  ];

  CLACK_SET = "Ashlynn,Terry Davis,Dennis Ritchie";
}

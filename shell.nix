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
    # Rust
    cargo
    cargo-watch
    rls
    rustc
    rustfmt

    # system dependencies
    openssl
    pkg-config

    # kubernetes deployment
    dhall
    dhall-yaml

    # dependency manager
    niv

    # tools
    ispell
  ];

  SITE_PREFIX = "devel.";
  CLACK_SET = "Ashlynn,Terry Davis,Dennis Ritchie";
  RUST_LOG = "debug";
  RUST_BACKTRACE = "1";
  GITHUB_SHA = "devel";
}

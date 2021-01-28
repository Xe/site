let
  sources = import ./nix/sources.nix;
  pkgs =
    import sources.nixpkgs { overlays = [ (import sources.nixpkgs-mozilla) ]; };
  dhallpkgs = import sources.easy-dhall-nix { inherit pkgs; };
  dhall-yaml = dhallpkgs.dhall-yaml-simple;
  dhall = dhallpkgs.dhall-simple;
  xepkgs = import sources.xepkgs { inherit pkgs; };
  rust = pkgs.callPackage ./nix/rust.nix { };
in with pkgs;
with xepkgs;
mkShell {
  buildInputs = [
    # Rust
    rust
    cargo-watch

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
  RUST_SRC_PATH =
    "${pkgs.latest.rustChannels.nightly.rust-src}/lib/rustlib/src/rust/library";
  GITHUB_SHA = "devel";
}

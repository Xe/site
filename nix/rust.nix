{ sources ? import ./sources.nix }:

let
  pkgs =
    import sources.nixpkgs { overlays = [ (import sources.nixpkgs-mozilla) ]; };
  channel = "nightly";
  date = "2020-11-25";
  targets = [ ];
  chan = pkgs.latest.rustChannels.stable.rust;
in chan

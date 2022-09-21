{ sources ? import ./sources.nix }:

let
  pkgs =
    import sources.nixpkgs { overlays = [ (import sources.nixpkgs-mozilla) ]; };
  channel = "nightly";
  date = "2022-09-15";
  targets = [ ];
  chan = pkgs.rustChannelOfTargets channel date targets;
in chan

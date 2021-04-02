{ sources ? import ./sources.nix }:

let
  pkgs =
    import sources.nixpkgs { overlays = [ (import sources.nixpkgs-mozilla) ]; };
  channel = "nightly";
  date = "2021-03-25";
  targets = [ ];
  chan = pkgs.rustChannelOfTargets channel date targets;
in chan

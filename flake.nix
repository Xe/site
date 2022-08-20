{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    naersk.url = "github:nix-community/naersk";
  };

  outputs = { self, nixpkgs, flake-utils, naersk }:
    flake-utils.lib.eachSystem [ "x86_64-linux" "aarch64-linux" ] (system:
      let
        pkgs = import nixpkgs { inherit system; };
        naersk-lib = naersk.lib."${system}";
        src = ./.;
      in rec {
        packages = rec {
          bin = naersk-lib.buildPackage {
            pname = "xesite-bin";
            root = src;
            buildInputs = with pkgs; [ pkg-config openssl git ];
          };

          config = pkgs.stdenv.mkDerivation {
            pname = "xesite-config";
            inherit (bin) version;
            inherit src;
            buildInputs = with pkgs; [ dhall ];

            phases = "installPhase";

            installPhase = ''
              cd $src
              mkdir -p $out
              dhall resolve < $src/config.dhall >> $out/config.dhall
            '';
          };

          static = pkgs.stdenv.mkDerivation {
            pname = "xesite-static";
            inherit (bin) version;
            inherit src;

            phases = "installPhase";

            installPhase = ''
              mkdir -p $out
              cp -vrf $src/static $out
              cp -vrf $src/css $out
            '';
          };

          posts = pkgs.stdenv.mkDerivation {
            pname = "xesite-posts";
            inherit (bin) version;
            inherit src;

            phases = "installPhase";

            installPhase = ''
              mkdir -p $out
              cp -vrf $src/blog $out
              cp -vrf $src/gallery $out
              cp -vrf $src/talks $out
            '';
          };

          default = pkgs.symlinkJoin {
            name = "xesite-${bin.version}";
            paths = [ config posts static bin ];
          };
        };

        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Rust
            rustc
            cargo
            rust-analyzer
            cargo-watch
            rustfmt

            # system dependencies
            openssl
            pkg-config

            # kubernetes deployment
            dhall
            dhall-json

            # dependency manager
            niv

            # tools
            ispell
            pandoc
          ];

          SITE_PREFIX = "devel.";
          CLACK_SET = "Ashlynn,Terry Davis,Dennis Ritchie";
          RUST_LOG = "debug";
          RUST_BACKTRACE = "1";
          GITHUB_SHA = "devel";
        };

        nixosModules.bot = { config, lib, ... }:
          with lib;
          let cfg = config.xeserv.services.xesite;
          in {
            options.within.services.xesite = {
              enable = mkEnableOption "Activates my personal website";
              useACME = mkEnableOption "Enables ACME for cert stuff";

              port = mkOption {
                type = types.port;
                default = 32837;
                example = 9001;
                description =
                  "The port number xesite should listen on for HTTP traffic";
              };

              domain = mkOption {
                type = types.str;
                default = "xesite.akua";
                example = "xeiaso.net";
                description =
                  "The domain name that nginx should check against for HTTP hostnames";
              };

              sockPath = mkOption rec {
                type = types.str;
                default = "/srv/within/run/xesite.sock";
                example = default;
                description =
                  "The unix domain socket that xesite should listen on";
              };
            };

            config = mkIf cfg.enable {
              users.users.xesite = {
                createHome = true;
                description = "github.com/Xe/site";
                isSystemUser = true;
                group = "within";
                home = "/srv/within/xesite";
                extraGroups = [ "keys" ];
              };

              systemd.services.xesite = {
                wantedBy = [ "multi-user.target" ];

                serviceConfig = {
                  User = "xesite";
                  Group = "within";
                  Restart = "on-failure";
                  WorkingDirectory = "/srv/within/xesite";
                  RestartSec = "30s";
                  Type = "notify";

                  # Security
                  CapabilityBoundingSet = "";
                  DeviceAllow = [ ];
                  NoNewPrivileges = "true";
                  ProtectControlGroups = "true";
                  ProtectClock = "true";
                  PrivateDevices = "true";
                  PrivateUsers = "true";
                  ProtectHome = "true";
                  ProtectHostname = "true";
                  ProtectKernelLogs = "true";
                  ProtectKernelModules = "true";
                  ProtectKernelTunables = "true";
                  ProtectSystem = "true";
                  ProtectProc = "invisible";
                  RemoveIPC = "true";
                  RestrictSUIDSGID = "true";
                  RestrictRealtime = "true";
                  SystemCallArchitectures = "native";
                  SystemCallFilter = [
                    "~@reboot"
                    "~@module"
                    "~@mount"
                    "~@swap"
                    "~@resources"
                    "~@cpu-emulation"
                    "~@obsolete"
                    "~@debug"
                    "~@privileged"
                  ];
                  UMask = "007";
                };

                script = let site = packages.default;
                in ''
                  export SOCKPATH=${cfg.sockPath}
                  export DOMAIN=${toString cfg.domain}
                  cd ${site}
                  exec ${site}/bin/xesite
                '';
              };

              services.nginx.virtualHosts."xesite" = {
                serverName = "${cfg.domain}";
                locations."/" = {
                  proxyPass = "http://unix:${toString cfg.sockPath}";
                  proxyWebsockets = true;
                };
                forceSSL = cfg.useACME;
                useACMEHost = "xeiaso.net";
                extraConfig = ''
                  access_log /var/log/nginx/xesite.access.log;
                '';
              };
            };
          };
      });
}

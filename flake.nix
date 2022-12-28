{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
    flake-utils.url = "github:numtide/flake-utils";
    naersk = {
      url = "github:nix-community/naersk";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    deno2nix = {
      url = "github:Xe/deno2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
  };

  outputs = { self, nixpkgs, flake-utils, naersk, deno2nix, ... }:
    flake-utils.lib.eachSystem [ "x86_64-linux" "aarch64-linux" ] (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ deno2nix.overlays.default ];
        };
        naersk-lib = naersk.lib."${system}";
        src = ./.;

        tex = with pkgs;
          texlive.combine { inherit (texlive) scheme-medium bitter titlesec; };
      in rec {
        packages = rec {
          bin = naersk-lib.buildPackage {
            pname = "xesite-bin";
            root = src;
            buildInputs = with pkgs; [
              pkg-config
              openssl
              git
              deno
              nodePackages.uglify-js
            ];
          };

          config = pkgs.stdenv.mkDerivation {
            pname = "xesite-config";
            inherit (bin) version;
            inherit src;
            buildInputs = with pkgs; [ dhall dhallPackages.Prelude ];

            phases = "installPhase";

            installPhase = ''
              mkdir -p $out
              cp -rf ${pkgs.dhallPackages.Prelude}/.cache .cache
              chmod -R u+w .cache
              export XDG_CACHE_HOME=.cache
              export DHALL_PRELUDE=${pkgs.dhallPackages.Prelude}/binary.dhall;
              dhall resolve --file $src/config.dhall >> $out/config.dhall
            '';
          };

          resumePDF = pkgs.stdenv.mkDerivation {
            pname = "xesite-resume-pdf";
            inherit (bin) version;
            inherit src;
            buildInputs = with pkgs; [ dhall dhallPackages.Prelude tex pandoc ];

            phases = "installPhase";

            installPhase = ''
              mkdir -p $out/static/resume
              cp -rf ${pkgs.dhallPackages.Prelude}/.cache .cache
              chmod -R u+w .cache
              export XDG_CACHE_HOME=.cache
              export DHALL_PRELUDE=${pkgs.dhallPackages.Prelude}/binary.dhall;

              ln -s $src/dhall/latex/resume.cls
              dhall text --file $src/dhall/latex/resume.dhall > resume.tex

              xelatex ./resume.tex
              cp resume.pdf $out/static/resume/resume.pdf
            '';
          };

          frontend = let
            build = { entrypoint, name ? entrypoint, minify ? true }: pkgs.deno2nix.mkBundled {
              pname = "xesite-frontend-${name}";
              inherit (bin) version;

              src = ./src/frontend;
              lockfile = ./src/frontend/deno.lock;

              output = "${entrypoint}.js";
              outPath = "static/js";
              entrypoint = "./${entrypoint}.tsx";
              importMap = "./import_map.json";
              inherit minify;
            };
            share-button = build {
              entrypoint = "mastodon_share_button";
            };
            wasiterm = build {
              entrypoint = "wasiterm";
            };
          in pkgs.symlinkJoin {
            name = "xesite-frontend-${bin.version}";
            paths = [ share-button wasiterm ];
          };

          static = pkgs.stdenv.mkDerivation {
            pname = "xesite-static";
            inherit (bin) version;
            inherit src;

            phases = "installPhase";

            installPhase = ''
              mkdir -p $out
              cp -vrf $src/data $out
              cp -vrf $src/static $out
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
            paths = [ config posts static bin frontend resumePDF ];
          };

          docker = pkgs.dockerTools.buildLayeredImage {
            name = "xena/xesite";
            tag = bin.version;
            contents = [ default ];
            config = {
              Cmd = [ "${bin}/bin/xesite" ];
              WorkdingDir = "${default}";
            };
          };
        };

        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Rust
            rustc
            cargo
            rust-analyzer
            cargo-watch
            cargo-license
            rustfmt

            # system dependencies
            openssl
            pkg-config

            # dhall
            dhall
            dhall-json
            dhall-lsp-server
            tex
            pandoc

            # frontend
            deno
            nodePackages.uglify-js

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
          DHALL_PRELUDE = "${pkgs.dhallPackages.Prelude}";
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

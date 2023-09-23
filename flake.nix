{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";

    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };

    naersk = {
      url = "github:nix-community/naersk";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    deno2nix = {
      url = "github:Xe/deno2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };

    # Explicitly pulling from that version of nixpkgs to avoid font duplication.
    iosevka.url = "github:Xe/iosevka";

    typst.url = "github:typst/typst";
    typst.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs =
    { self, nixpkgs, flake-utils, naersk, deno2nix, iosevka, typst, ... }:
    flake-utils.lib.eachSystem [ "x86_64-linux" "aarch64-linux" "aarch64-darwin" ] (system:
      let
        graft = pkgs: pkg: pkg.override {
          buildGoModule = pkgs.buildGo121Module;
        };
        pkgs = import nixpkgs {
          inherit system;
          overlays = [
            deno2nix.overlays.default
            typst.overlays.default
            (final: prev: {
              go = prev.go_1_21;
              go-tools = graft prev prev.go-tools;
              gotools = graft prev prev.gotools;
              gopls = graft prev prev.gopls;
            })
          ];
        };
        naersk-lib = naersk.lib."${system}";
        src = ./.;
        lib = pkgs.lib;

        tex = with pkgs;
          texlive.combine { inherit (texlive) scheme-medium bitter titlesec; };

        fontsConf = pkgs.symlinkJoin {
            name = "typst-fonts";
            paths = [ "${self.packages.${system}.iosevka}/static/css/iosevka" ];
          };

        typstWithIosevka = pkgs.writeShellApplication {
          name = "typst";
          text = ''
            ${pkgs.typst-dev}/bin/typst \
            --font-path ${fontsConf} \
            "$@"
          '';
          runtimeInputs = [ ];
        };
      in rec {
        packages = rec {
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
            buildInputs = with pkgs; [
              dhall-json
              dhallPackages.Prelude
              tex
              pandoc
              typst-dev
            ];

            phases = "installPhase";

            installPhase = ''
              mkdir -p $out/static/resume

              cp -rf ${pkgs.dhallPackages.Prelude}/.cache .cache
              chmod -R u+w .cache
              export XDG_CACHE_HOME=.cache
              export DHALL_PRELUDE=${pkgs.dhallPackages.Prelude}/binary.dhall;

              mkdir -p icons
              cp -vrf $src/dhall/resume/* .
              dhall-to-json --file $src/dhall/resume.dhall --output resume.json

              typst compile --font-path ${fontsConf} resume.typ $out/static/resume/resume.pdf
            '';
          };

          default = pkgs.symlinkJoin {
            name = "xesite-${bin.version}";
            paths = [ config resumePDF ];
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
            # Go
            go
            go-tools
            gotools
            gopls

            # dhall
            dhall
            dhall-json
            tex
            pandoc
            typstWithIosevka
            pagefind

            # frontend
            deno
            nodePackages.uglify-js
            esbuild
            zig
            nodejs

            # tools
            ispell
            pandoc
            python311Packages.fonttools
          ];

          SITE_PREFIX = "devel.";
          CLACK_SET = "Ashlynn,Terry Davis,Dennis Ritchie";
          ESBUILD_BINARY_PATH = "${pkgs.esbuild}/bin/esbuild";
          RUST_LOG = "debug";
          RUST_BACKTRACE = "1";
          GITHUB_SHA = "devel";
          DHALL_PRELUDE = "${pkgs.dhallPackages.Prelude}";
        };
      }) // {
        nixosModules.default = import ./nix/xesite.nix self;
      };
}

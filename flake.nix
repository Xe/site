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

    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };

    # Explicitly pulling from that version of nixpkgs to avoid font duplication.
    iosevka.url = "github:Xe/iosevka";

    typst.url = "github:typst/typst";
    typst.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs =
    { self, nixpkgs, flake-utils, naersk, deno2nix, iosevka, typst, gomod2nix, ... }:
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
            gomod2nix.overlays.default
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
            compile \
            --font-path ${fontsConf} \
            "$@"
          '';
          runtimeInputs = [ ];
        };

        # Generate a user-friendly version number.
      version = builtins.substring 0 8 self.lastModifiedDate;
      in rec {
        packages = rec {
          bin = pkgs.buildGoApplication {
            pname = "xesite_v4";
            inherit version;
            src = ./.;
            modules = ./gomod2nix.toml;
            subPackages = [ "xesite" ];
          };

          docker = pkgs.dockerTools.buildLayeredImage {
            name = "xena/xesite";
            tag = version;
            contents = with pkgs; [ ca-certificates typstWithIosevka dhall-json deno ];
            config = {
              Cmd = [ "${bin}/bin/xesite" ];
              Env = [
                "TMPDIR=/data"
              ];
              Volumes = {
                "/data" = {};
              };
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
            #typstWithIosevka
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

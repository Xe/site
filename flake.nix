{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";

    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
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

    alpineLinux = {
      flake = false;
      url =
        "file+https://cdn.xeiaso.net/file/christine-static/hack/alpine-amd64-3.19.0-1.tar";
    };
  };

  outputs = { self, nixpkgs, flake-utils, deno2nix, iosevka, typst, gomod2nix
    , alpineLinux, ... }:
    flake-utils.lib.eachSystem [
      "x86_64-linux"
      "aarch64-linux"
      "aarch64-darwin"
    ] (system:
      let
        graft = pkgs: pkg:
          pkg.override { buildGoModule = pkgs.buildGo122Module; };
        pkgs = import nixpkgs {
          inherit system;
          overlays = [
            deno2nix.overlays.default
            typst.overlays.default
            (final: prev: {
              go = prev.go_1_22;
              go-tools = graft prev prev.go-tools;
              gotools = graft prev prev.gotools;
              gopls = graft prev prev.gopls;
            })
            gomod2nix.overlays.default
          ];
        };
        src = ./.;
        lib = pkgs.lib;

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
            subPackages = [ "cmd/xesite" ];
          };

          patreon-bin = pkgs.buildGoApplication {
            pname = "patreon-saasproxy";
            inherit version;
            src = ./.;
            modules = ./gomod2nix.toml;
            subPackages = [ "cmd/patreon-saasproxy" ];
          };

          iosevka = pkgs.stdenvNoCC.mkDerivation {
            name = "xesite-iosevka";
            buildInputs = with pkgs; [
              python311Packages.brotli
              python311Packages.fonttools
            ];
            dontUnpack = true;
            buildPhase = ''
              mkdir -p out
              ${pkgs.unzip}/bin/unzip ${
                self.inputs.iosevka.packages.${system}.default
              }/ttf.zip
              for ttf in ttf/*.ttf; do
                cp $ttf out
                name=`basename -s .ttf $ttf`
                pyftsubset \
                    $ttf \
                    --output-file=out/"$name".woff2 \
                    --flavor=woff2 \
                    --layout-features=* \
                    --no-hinting \
                    --desubroutinize \
                    --unicodes="U+0000-0170,U+00D7,U+00F7,U+2000-206F,U+2074,U+20AC,U+2122,U+2190-21BB,U+2212,U+2215,U+F8FF,U+FEFF,U+FFFD,U+00E8"
              done
            '';
            installPhase = ''
              mkdir -p $out/static/css/iosevka
              cp out/* $out/static/css/iosevka
            '';
          };

          docker = pkgs.dockerTools.buildLayeredImage {
            name = "ghcr.io/xe/site/bin";
            tag = "latest";
            fromImage = alpineLinux;
            contents = with pkgs; [ cacert typst-dev dhall-json deno git ];
            config = {
              Cmd = [ "${bin}/bin/xesite" "--data-dir=/data" ];
              Env = [
                "HOME=/data"
                "DHALL_PRELUDE=${pkgs.dhallPackages.Prelude}"
                "TYPST_FONT_PATHS=${fontsConf}"
              ];
              Volumes."/data" = { };
            };
          };

          patreon-docker = pkgs.dockerTools.buildLayeredImage {
            name = "ghcr.io/xe/site/patreon";
            tag = "latest";
            contents = with pkgs; [ cacert ];
            config = {
              Cmd = [ "${patreon-bin}/bin/patreon-saasproxy" ];
              Env = [ "HOME=/data" ];
              Volumes."/data" = { };
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
            gomod2nix.packages.${system}.default

            # dhall
            dhall
            dhall-json
            typst-dev
            pagefind

            # frontend
            deno
            nodePackages.uglify-js
            esbuild
            zig
            nodejs

            protobuf
            protoc-gen-go
            protoc-gen-twirp

            jq
            jo

            earthly

            # tools
            ispell
            pandoc
            python311Packages.fonttools
          ];

          DHALL_PRELUDE = "${pkgs.dhallPackages.Prelude}";
          TYPST_FONT_PATHS = "${fontsConf}";
          FLY_REGION = "dev";
        };
      });
}

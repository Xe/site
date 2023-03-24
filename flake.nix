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
  };

  outputs = { self, nixpkgs, flake-utils, naersk, deno2nix, iosevka, ... }:
    flake-utils.lib.eachSystem [ "x86_64-linux" "aarch64-linux" ] (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ deno2nix.overlays.default ];
        };
        naersk-lib = naersk.lib."${system}";
        src = ./.;
        lib = pkgs.lib;

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
            build = { entrypoint, name ? entrypoint, minify ? true }:
              pkgs.deno2nix.mkBundled {
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
            share-button = build { entrypoint = "mastodon_share_button"; };
            wasiterm = build { entrypoint = "wasiterm"; };
          in pkgs.symlinkJoin {
            name = "xesite-frontend-${bin.version}";
            paths = [ share-button wasiterm ];
          };

          iosevka = pkgs.stdenvNoCC.mkDerivation {
            name = "xesite-iosevka";
            buildInputs = with pkgs; [ python311Packages.brotli python311Packages.fonttools ];
            dontUnpack = true;
            buildPhase =
              ''
                mkdir -p out
                ${pkgs.unzip}/bin/unzip ${self.inputs.iosevka.packages.${system}.default}/ttf.zip
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
                      --unicodes="U+0000-00A0,U+00A2-00A9,U+00AC-00AE,U+00B0-00B7,U+00B9-00BA,U+00BC-00BE,U+00D7,U+00F7,U+2000-206F,U+2074,U+20AC,U+2122,U+2190-21BB,U+2212,U+2215,U+F8FF,U+FEFF,U+FFFD"
                done

              '';
            installPhase = ''
              mkdir -p $out/static/css/iosevka
              cp out/* $out/static/css/iosevka
            '';
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
            paths = [ config posts static bin frontend resumePDF iosevka ];
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
            python311Packages.fonttools
          ];

          SITE_PREFIX = "devel.";
          CLACK_SET = "Ashlynn,Terry Davis,Dennis Ritchie";
          RUST_LOG = "debug";
          RUST_BACKTRACE = "1";
          GITHUB_SHA = "devel";
          DHALL_PRELUDE = "${pkgs.dhallPackages.Prelude}";
        };
      }) // {
        nixosModules.default = import ./nix/xesite.nix self;
      };
}

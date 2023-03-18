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
            name = "custom-iosevka";
            dontUnpack = true;
            buildPhase =
              let
                metric-override = {
                  cap = 790;
                  ascender = 790;
                  xHeight = 570;
                };
                iosevka-term = pkgs.iosevka.override {
                  set = "curly";
                  privateBuildPlan = {
                    family = "Iosevka Term Iaso";
                    spacing = "term";
                    serifs = "sans";
                    no-ligation = false;
                    ligations = {
                      "inherit" = "default-calt";
                    };
                    no-cv-ss = true;
                    variants = {
                      inherits = "ss01";
                      design = {
                        tilde = "low";
                        number-sign = "slanted-tall";
                        at = "fourfold-solid-inner-tall";
                      };
                    };
                    slopes.upright = {
                      angle = 0;
                      shape = "upright";
                      menu = "upright";
                      css = "normal";
                    };
                    weights.regular = {
                      shape = 400;
                      menu = 400;
                      css = 400;
                    };
                    widths.normal = {
                      shape = 600;
                      menu = 7;
                      css = "expanded";
                    };
                    inherit metric-override;
                  };
                };
                iosevka-aile = pkgs.iosevka.override {
                  set = "aile";
                  privateBuildPlan = {
                    family = "Iosevka Aile Iaso";
                    spacing = "quasi-proportional-extension-only";
                    no-ligation = true;
                    no-cv-ss = true;
                    variants = {
                      inherits = "ss01";
                      design = {
                        tilde = "low";
                        number-sign = "slanted-tall";
                        at = "fourfold-solid-inner-tall";
                      };
                    };
                    slopes = {
                      upright = {
                        angle = 0;
                        shape = "upright";
                        menu = "upright";
                        css = "normal";
                      };
                      italic = {
                        angle = 9.4;
                        shape = "italic";
                        menu = "italic";
                        css = "italic";
                      };
                    };
                    weights.regular = {
                      shape = 400;
                      menu = 400;
                      css = 400;
                    };
                    widths.normal = {
                      shape = 550;
                      menu = 7;
                      css = "expanded";
                    };
                    inherit metric-override;
                  };
                };
                iosevka-etoile = pkgs.iosevka.override {
                  set = "etoile";
                  privateBuildPlan = {
                    family = "Iosevka Etoile Iaso";
                    spacing = "quasi-proportional";
                    serifs = "slab";
                    no-ligation = true;
                    no-cv-ss = true;
                    variants = {
                      inherits = "ss01";
                      design = {
                        capital-w = "straight-flat-top";
                        f = "flat-hook-serifed";
                        j = "flat-hook-serifed";
                        t = "flat-hook";
                        capital-t = "serifed"; # not part of original Iosevka Aile
                        w = "straight-flat-top";
                        #capital-g = "toothless-rounded-serifless-hooked";
                        r = "corner-hooked";

                        tilde = "low";
                        number-sign = "slanted-tall";
                        at = "fourfold-solid-inner-tall";
                      };
                      italic = {
                        f = "flat-hook-tailed";
                      };
                    };
                    slopes = {
                      upright = {
                        angle = 0;
                        shape = "upright";
                        menu = "upright";
                        css = "normal";
                      };
                      italic = {
                        angle = 9.4;
                        shape = "italic";
                        menu = "italic";
                        css = "italic";
                      };
                    };
                    weights.regular = {
                      shape = 400;
                      menu = 400;
                      css = 400;
                    };
                    widths.normal = {
                      shape = 600;
                      menu = 7;
                      css = "expanded";
                    };
                    inherit metric-override;
                  };
                };
              in
              ''
                mkdir -p ttf
                for ttf in ${iosevka-term}/share/fonts/truetype/*.ttf ${iosevka-aile}/share/fonts/truetype/*.ttf ${iosevka-etoile}/share/fonts/truetype/*.ttf; do
                  cp $ttf .
                  ${pkgs.woff2}/bin/woff2_compress *.ttf
                  mv *.ttf ttf
                done
              '';
            installPhase = ''
              mkdir -p $out/static/css/iosevka
              cp *.woff2 $out/static/css/iosevka
              cp ttf/*.ttf $out/static/css/iosevka
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

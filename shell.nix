let
  sources = import ./nix/sources.nix;
  pkgs = import sources.nixpkgs { };
  niv = (import sources.niv { }).niv;
  dhall-yaml =
    (import sources.easy-dhall-nix { inherit pkgs; }).dhall-yaml-simple;
  xepkgs = import sources.xepkgs { inherit pkgs; };
  vgo2nix = import sources.vgo2nix { inherit pkgs; };
in pkgs.mkShell {
  buildInputs = [ pkgs.go xepkgs.gopls dhall-yaml niv vgo2nix ];
}

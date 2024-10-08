{
  description = "A helper utility for my desktop";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    flake-parts = {
      url = "github:hercules-ci/flake-parts";
    };

    devenv = {
      url = "github:cachix/devenv";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{ flake-parts, nixpkgs, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [ inputs.devenv.flakeModule ];
      systems = nixpkgs.lib.systems.flakeExposed;

      perSystem =
        {
          config,
          self',
          inputs',
          pkgs,
          system,
          ...
        }:
        {
          devenv.shells.default = {
            packages = with pkgs; [ just gopls hyprland ];
            languages.nix.enable = true;
            languages.go.enable = true;
            languages.go.enableHardeningWorkaround = true;
          };
          packages = rec {
            dfh = pkgs.callPackage ./package.nix { version = "0.1"; };
            default = dfh;
          };
        };
    };
}
